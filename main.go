package main

import log "github.com/Sirupsen/logrus"

var (
	JobQueue      chan Job
	pages         chan Page
	processed     chan Job
	lineEnqueued  chan bool
	lineProcessed chan bool
	allDone       chan bool

	source  Source
	service *Service
)

type Page struct {
	count int
	page  int
	last  bool
}

type Job struct {
	id   int
	key  string
	page int
}

func init() {
	logger(err_log, "errors.log")
	logger(info_log, "info.log")
	options = getOptions()
}

func main() {
	service = NewService()
	if options.Reader == "file" {
		source = NewFileSource(options.FilePath)
		lineEnqueued = make(chan bool, options.Queue)
		lineProcessed = make(chan bool, options.Queue)
	} else {
		source = NewS3Source()
		pages = make(chan Page, 1000)
		processed = make(chan Job, options.Queue)
	}

	JobQueue = make(chan Job, options.Queue)
	allDone = make(chan bool, 1)

	dispatcher := NewDispatcher(options.Workers)
	dispatcher.Run()

	source.Run()
}

// watches each page for added and processed objects, determines when it's done and logs finished
func watch() {
	lines := 0
	done := false
	p := make(map[int]int) //page:count

	for {
		select {
		case page := <-pages:
			p[page.page] = page.count
			done = page.last
		case job := <-processed:
			p[job.page]--
			if p[job.page] == 0 {
				info_log.WithFields(log.Fields{"Key": job.key, "Page": job.page, "Done": true}).
					Info("Page ", job.page, " finished all jobs")
				delete(p, job.page)
			}
			if done && len(p) == 0 {
				allDone <- true
				return
			}
		case <-lineEnqueued:
			lines++
		case <-lineProcessed:
			lines--
			if lines == 0 {
				allDone <- true
				return
			}

		}
	}
}
