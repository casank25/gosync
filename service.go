package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Service struct {
	Sess   *session.Session
	Client *s3.S3
}

func NewService() *Service {
	sess, sess_err := session.NewSession(&aws.Config{
		Region:     aws.String(options.Region),
		MaxRetries: &options.Retries,
	},
	)
	if sess_err != nil {
		log.Fatal("Session Error: ", sess_err.Error())
	}

	client := s3.New(sess)

	return &Service{
		Sess:   sess,
		Client: client,
	}
}
