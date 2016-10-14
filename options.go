package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
)

var options Option

type Option struct {
	Reader      string
	FilePath    string
	Source      string
	Destination string
	Region      string
	Start       string
	Retries     int
	Workers     int
	Queue       int
}

func getOptions() Option {
	reader := flag.String("reader", "s3", "read from file or s3 pages")
	filePath := flag.String("file-path", "", "file path if reader set to file")
	source := flag.String("source", "", "source bucket")
	destination := flag.String("destination", "", "destination bucket")
	region := flag.String("region", "us-west-2", "awa bucket region")
	retries := flag.Int("retries", 5, "max retries")
	start := flag.String("start", "", "start from key")
	workers := flag.Int("workers", 100, "max number of workers")
	queue := flag.Int("queue", 1000, "Max Queue Size")

	flag.Parse()

	if *source == "" || *destination == "" {
		log.Fatal("source and destination arguments are required")
	}

	return Option{
		Reader:      *reader,
		FilePath:    *filePath,
		Source:      *source,
		Destination: *destination,
		Region:      *region,
		Start:       *start,
		Retries:     *retries,
		Workers:     *workers,
		Queue:       *queue,
	}
}
