package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func setup() (*session.Session, error) {
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}

func NewS3Service() (*s3.S3, error) {
	session, err := setup()
	if err != nil {
		return nil, err
	}
	s3Session := s3.New(session)

	return s3Session, nil
}
