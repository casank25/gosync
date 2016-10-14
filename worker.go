package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				w.process(job)
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func (w Worker) process(j Job) {
	copySource := options.Source + "/" + j.key
	params := &s3.CopyObjectInput{
		Bucket:     aws.String(options.Destination),
		CopySource: aws.String(copySource),
		Key:        aws.String(j.key),
	}

	fmt.Println("Copying Object: ", j.key)
	_, err := service.Client.CopyObject(params)

	switch options.Reader {
	case "s3":
		processed <- j
	case "file":
		lineProcessed <- true
	}

	if err != nil {
		err_log.WithFields(log.Fields{"Name": j.key}).
			Error("Copy Error:", err.Error())
	}
}
