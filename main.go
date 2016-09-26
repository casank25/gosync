package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var JobQueue chan Job
var pages chan Page
var processed chan Job
var allDone chan bool

var svc *s3.S3

type Page struct {
	count int
	page  int
	last  bool
}

type Job struct {
	id     int
	object *s3.Object
	page   int
}

func main() {
	logger(err_log, "errors.log")
	logger(info_log, "info.log")
	options = getOptions()

	sess, sess_err := session.NewSession(&aws.Config{
		Region:     aws.String(options.Region),
		MaxRetries: &options.Retries,
	},
	)
	if sess_err != nil {
		log.Fatal("Session Error: ", sess_err.Error())
	}

	JobQueue = make(chan Job, options.Queue)
	pages = make(chan Page, 1000)
	processed = make(chan Job, options.Queue)
	allDone = make(chan bool, 1)

	dispatcher := NewDispatcher(options.Workers)
	dispatcher.Run()

	run(sess)

}

func run(sess *session.Session) {
	svc = s3.New(sess)
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(options.Source),
	}
	if options.Start != "" {
		params.StartAfter = &options.Start
	}

	go watch()
	c := 0
	list_err := svc.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			enqueued := 0
			c++

			for _, obj := range page.Contents {
				enqueued++
				job := Job{id: enqueued, object: obj, page: c}
				JobQueue <- job
			}
			p := Page{count: enqueued, page: c, last: lastPage}
			pages <- p

			if lastPage {
				return false
			}

			return true
		},
	)

	if list_err != nil {
		log.Fatal("Could not list pages: ", list_err.Error())
	}

	<-allDone
}

// watches each page for added and processed objects, determines when it's done and logs finished
func watch() {
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
				info_log.WithFields(log.Fields{"Key": *job.object.Key, "Page": job.page, "Done": true}).
					Info("Page ", job.page, " finished all jobs")
				delete(p, job.page)
			}
			if done && len(p) == 0 {
				allDone <- true
				return
			}
		}
	}
}
