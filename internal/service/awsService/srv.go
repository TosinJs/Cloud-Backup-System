package awsService

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"mime/multipart"
)

type AWSService struct {
	s3 *s3.S3
}

func New(s3 *s3.S3) AWSService {
	return AWSService{
		s3: s3,
	}
}

func (a AWSService) UploadFile(file multipart.File, filename string) error {

	input := &s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String("cloud-backup-system"),
		Key:    aws.String(filename),
	}

	result, err := a.s3.PutObject(input)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

func (a AWSService) DeleteFile(filename string) error {
	res, err := a.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("cloud-backup-system"),
		Key:    aws.String(filename),
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(res)
	return nil
}

func (a AWSService) DeleteFolder(foldername string) error {
	result, err := a.listFilesInFolder(foldername)

	if err != nil {
		return err
	}

	if len(result.Contents) < 1 {
		return fmt.Errorf("this folder does not exist")
	}

	if len(result.CommonPrefixes) != 0 {
		return fmt.Errorf("delete sub directories before the parent directory")
	}

	deletedArray := []*s3.ObjectIdentifier{}

	for _, i := range result.Contents {
		deletedArray = append(deletedArray, &s3.ObjectIdentifier{Key: i.Key})
	}

	_, err := a.s3.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String("cloud-backup-system"),
		Delete: &s3.Delete{
			Objects: deletedArray,
		},
	})

	return nil
}

func (a AWSService) GetFile(filename string) ([]byte, error) {
	result, err := a.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("cloud-backup-system"),
		Key:    aws.String(filename),
	})

	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	return body, nil
}

func (a AWSService) ListFilesInFolder(foldername string) ([]string, error) {
	result, err := a.listFilesInFolder(foldername)

	if err != nil {
		return nil, err
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

func (a AWSService) listFilesInFolder(foldername string) (*s3.ListObjectsV2Output, error) {
	return a.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String("cloud-backup-system"),
		Prefix:    aws.String(foldername),
		Delimiter: aws.String("/"),
	})
}
