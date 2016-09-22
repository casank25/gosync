package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var JobQueue chan Job
var pages chan Page
var processed chan bool
var allDone chan bool
var interrupt chan os.Signal
var lastObject string

var svc *s3.S3

type Page struct {
	count int
	last  bool
}

type Job struct {
	id     int
	object *s3.Object
}

func main() {
	logger(err_log, "errors.log")
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
	processed = make(chan bool, options.Queue)
	allDone = make(chan bool, 1)
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go listen()

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

	go track()
	list_err := svc.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			enqueued := 0
			for _, obj := range page.Contents {
				enqueued++
				job := Job{id: enqueued, object: obj}
				JobQueue <- job
			}
			p := Page{count: enqueued, last: lastPage}
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

// tracks for added and processed objects and defines when all of them have been copied
func track() {
	done := false
	count := 0

	for {
		select {
		case page := <-pages:
			count = count + page.count
			done = page.last
		case <-processed:
			count--
			if done && count == 0 {
				allDone <- true
				return
			}
		}
	}
}

// listens for CTRL-C command from user and records the last copied object
func listen() {
	<-interrupt
	close(JobQueue)
	close(pages)
	close(processed)
	err_log.Error("Program interrupted by user. Last object copied: ", lastObject)
	os.Exit(1)
}
