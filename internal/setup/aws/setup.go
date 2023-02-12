package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func setup(
	AWS_ID, AWS_SECRET, AWS_TOKEN, AWS_REGION string,
) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(
			AWS_ID,
			AWS_SECRET,
			AWS_TOKEN,
		),
	})
}

func NewS3Service(
	AWS_ID, AWS_SECRET, AWS_TOKEN, AWS_REGION string,
) (*s3.S3, error) {
	session, err := setup(AWS_ID, AWS_SECRET, AWS_TOKEN, AWS_REGION)
	if err != nil {
		return nil, err
	}
	s3Session := s3.New(session)
	return s3Session, nil
}
