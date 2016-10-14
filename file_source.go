package main

import (
	"bufio"
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
)

type FileSource struct {
	file    *os.File
	scanner *bufio.Scanner
	path    string
}

func NewFileSource(path string) *FileSource {
	return &FileSource{
		path: path,
	}
}

func (f *FileSource) Run() {
	go watch()

	f.Open()
	for f.scanner.Scan() {
		var object Object
		b := f.scanner.Bytes()
		if err := json.Unmarshal(b, &object); err == nil {
			err_log.Error("Could not convert json:", err)
		}

		job := Job{key: object.Key}
		JobQueue <- job
		lineEnqueued <- true
	}

	f.Close()
	<-allDone
}

func (f *FileSource) Open() *FileSource {
	var err error
	f.file, err = os.Open(f.path)
	if err != nil {
		log.Fatal("Could not open the file", err)
	}
	f.scanner = bufio.NewScanner(f.file)
	return f
}

func (f *FileSource) Close() *FileSource {
	err := f.file.Close()
	if err != nil {
		log.Error("Could not close file", err)
	}

	return f
}
