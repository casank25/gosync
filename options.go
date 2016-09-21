package main

import "flag"

var options Option

type Option struct {
	Source      string
	Destination string
	Region      string
	Start       string
	Retries     int
	Workers     int
	Queue       int
}

func getOptions() Option {
	source := flag.String("source", "com.boro-group.bucket3", "source bucket")
	destination := flag.String("destination", "com.boro-group.bucket2", "destination bucket")
	region := flag.String("region", "us-west-2", "awa bucket region")
	retries := flag.Int("retries", 5, "max retries")
	start := flag.String("start", "", "start from key")
	workers := flag.Int("workers", 100, "max number of workers")
	queue := flag.Int("queue", 1000, "Max Queue Size")

	flag.Parse()

	return Option{
		Source:      *source,
		Destination: *destination,
		Region:      *region,
		Start:       *start,
		Retries:     *retries,
		Workers:     *workers,
		Queue:       *queue,
	}
}
