package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

var err_log = log.New()
var info_log = log.New()

func logger(logger *log.Logger, filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Could not create log file: ", err.Error())
	}
	logger.Out = file
}
