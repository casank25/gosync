package main

import (
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
	copySource := options.Source + "/" + *j.object.Key
	params := &s3.CopyObjectInput{
		Bucket:     aws.String(options.Destination),
		CopySource: aws.String(copySource),
		Key:        j.object.Key,
	}

	_, err := svc.CopyObject(params)
	processed <- true
	if err != nil {
		err_log.WithFields(log.Fields{"Size": *j.object.Size, "Name": *j.object.Key}).
			Error("Copy Error:", err.Error())
	}
}
