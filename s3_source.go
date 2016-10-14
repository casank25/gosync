package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Source struct {
	sess *session.Session
}

func NewS3Source() *S3Source {
	return &S3Source{}
}

func (s *S3Source) Run() {
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(options.Source),
	}
	if options.Start != "" {
		params.StartAfter = &options.Start
	}

	go watch()
	c := 0
	list_err := service.Client.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			enqueued := 0
			c++

			for _, obj := range page.Contents {
				enqueued++
				job := Job{id: enqueued, key: *obj.Key, page: c}
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
