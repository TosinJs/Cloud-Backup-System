package awsService

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"mime/multipart"
	"tosinjs/cloud-backup/internal/entity/errorEntity"
)

type awsService struct {
	s3 *s3.S3
}

type AWSService interface {
	UploadFile(file multipart.File, filename string) *errorEntity.ServiceError
	DeleteFile(filename string) *errorEntity.ServiceError
	DeleteFolder(foldername string) *errorEntity.ServiceError
	GetFile(filename string) ([]byte, *errorEntity.ServiceError)
	ListFilesInFolder(foldername string) ([]string, *errorEntity.ServiceError)
}

func New(s3 *s3.S3) AWSService {
	return awsService{
		s3: s3,
	}
}

func (a awsService) UploadFile(file multipart.File, filename string) *errorEntity.ServiceError {

	input := &s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String("cloud-backup-system"),
		Key:    aws.String(filename),
	}

	_, err := a.s3.PutObject(input)
	if err != nil {
		return errorEntity.InternalServerError(err)
	}
	return nil
}

func (a awsService) DeleteFile(filename string) *errorEntity.ServiceError {
	_, err := a.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("cloud-backup-system"),
		Key:    aws.String(filename),
	})

	if err != nil {
		return errorEntity.InternalServerError(err)
	}

	return nil
}

func (a awsService) DeleteFolder(foldername string) *errorEntity.ServiceError {
	result, err := a.listFilesInFolder(foldername)

	if err != nil {
		return errorEntity.InternalServerError(err)
	}

	if len(result.Contents) < 1 {
		return errorEntity.NotFoundError("Folder / Path Not Found", errors.New("object key not found"))
	}

	if len(result.CommonPrefixes) != 0 {
		return errorEntity.BadRequestError("delete sub directories before the parent directory", errors.New("invalid delete operation"))
	}

	deletedArray := []*s3.ObjectIdentifier{}

	for _, i := range result.Contents {
		deletedArray = append(deletedArray, &s3.ObjectIdentifier{Key: i.Key})
	}

	_, err = a.s3.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String("cloud-backup-system"),
		Delete: &s3.Delete{
			Objects: deletedArray,
		},
	})

	if err != nil {
		return errorEntity.InternalServerError(err)
	}

	return nil
}

func (a awsService) GetFile(filename string) ([]byte, *errorEntity.ServiceError) {
	result, err := a.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("cloud-backup-system"),
		Key:    aws.String(filename),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return []byte{}, errorEntity.NotFoundError("Invalid File Path", aerr)
			default:
				return []byte{}, errorEntity.InternalServerError(aerr)
			}
		} else {
			return []byte{}, errorEntity.InternalServerError(err)
		}
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return []byte{}, errorEntity.InternalServerError(err)
	}

	return body, nil
}

func (a awsService) ListFilesInFolder(foldername string) ([]string, *errorEntity.ServiceError) {
	result, err := a.listFilesInFolder(foldername)

	if err != nil {
		return nil, errorEntity.InternalServerError(err)
	}

	fileStructure := []string{}

	for _, i := range result.CommonPrefixes {
		fileStructure = append(fileStructure, *i.Prefix)
	}

	for _, j := range result.Contents {
		fileStructure = append(fileStructure, *j.Key)
	}

	return fileStructure, nil
}

func (a awsService) listFilesInFolder(foldername string) (*s3.ListObjectsV2Output, error) {
	return a.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String("cloud-backup-system"),
		Prefix:    aws.String(foldername),
		Delimiter: aws.String("/"),
	})
}
