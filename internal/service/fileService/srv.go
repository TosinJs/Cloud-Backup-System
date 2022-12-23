package fileService

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"tosinjs/cloud-backup/internal/service/awsService"
)

type FileService struct {
	awsSVC awsService.AWSService
}

func New(awsSVC awsService.AWSService) FileService {
	return FileService{
		awsSVC: awsSVC,
	}
}

func (f FileService) UploadFile(hFile *multipart.FileHeader, folderName string) error {
	filename := hFile.Filename
	file, err := hFile.Open()
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	defer file.Close()

	if folderName != "" {
		filename = fmt.Sprintf("%s/%s/%s", "username", folderName, filename)
	} else {
		filename = fmt.Sprintf("%s/%s", "username", filename)
	}

	err = f.awsSVC.UploadFile(file, filename)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return err
	}

	return nil
}

func (f *FileService) DeleteFile(folderName, filename string) error {

	fileName := fmt.Sprintf("%s/", "username")

	if folderName != "" {
		fileName = fmt.Sprintf("%s%s/", fileName, folderName)
	}

	if fileName != "" {
		fileName = fmt.Sprintf("%s%s", fileName, filename)
	}

	fmt.Println(folderName, filename)
	fmt.Println(fileName)

	err := f.awsSVC.DeleteFile(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (f *FileService) GetFile(folderName, filename string) ([]byte, error) {
	fileName := fmt.Sprintf("%s/", "username")

	if folderName != "" {
		fileName = fmt.Sprintf("%s%s/", fileName, folderName)
	}

	if fileName != "" {
		fileName = fmt.Sprintf("%s%s", fileName, filename)
	}

	fmt.Println(folderName, filename)
	fmt.Println(fileName)

	data, err := f.awsSVC.GetFile(fileName)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	return data, nil
}

func (f *FileService) ListFilesInFolder(folderName string) ([]string, error) {
	fileName := fmt.Sprintf("%s/", "username")

	if folderName != "" {
		fileName = fmt.Sprintf("%s%s/", fileName, folderName)
	}

	fileStructure, err := f.awsSVC.ListFilesInFolder(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return fileStructure, nil
}

func (f *FileService) DeleteFolder(folderName string) error {
	fileName := fmt.Sprintf("%s/", "username")

	if folderName != "" {
		fileName = fmt.Sprintf("%s%s/", fileName, folderName)
	}

	err := f.awsSVC.DeleteFolder(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
