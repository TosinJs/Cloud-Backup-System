package fileService

import (
	"fmt"
	"mime/multipart"

	"tosinjs/cloud-backup/internal/entity/errorEntity"
	"tosinjs/cloud-backup/internal/service/awsService"
)

type fileService struct {
	awsSVC awsService.AWSService
}

type FileService interface {
	UploadFile(hFile *multipart.FileHeader, username, folderName string) *errorEntity.ServiceError
	DeleteFile(username, folderName, filename string) *errorEntity.ServiceError
	GetFile(username, folderName, filename string) ([]byte, *errorEntity.ServiceError)
	DeleteFolder(username, folderName string) *errorEntity.ServiceError
	ListFilesInFolder(username, folderName string) ([]string, *errorEntity.ServiceError)
}

func New(awsSVC awsService.AWSService) FileService {
	return fileService{
		awsSVC: awsSVC,
	}
}

func (f fileService) UploadFile(hFile *multipart.FileHeader, username, folderName string) *errorEntity.ServiceError {
	filename := hFile.Filename
	file, err := hFile.Open()
	if err != nil {
		return errorEntity.InternalServerError(err)
	}

	defer file.Close()

	if folderName != "" {
		filename = fmt.Sprintf("%s/%s/%s", username, folderName, filename)
	} else {
		filename = fmt.Sprintf("%s/%s", username, filename)
	}

	return f.awsSVC.UploadFile(file, filename)
}

func (f fileService) DeleteFile(username, folderName, filename string) *errorEntity.ServiceError {

	fileName := fmt.Sprintf("%s/", username)

	if folderName != "" {
		fileName = fmt.Sprintf("%s%s/", fileName, folderName)
	}

	if fileName != "" {
		fileName = fmt.Sprintf("%s%s", fileName, filename)
	}

	return f.awsSVC.DeleteFile(fileName)
}

func (f fileService) GetFile(username, folderName, filename string) ([]byte, *errorEntity.ServiceError) {
	fileName := fmt.Sprintf("%s/", username)

	if folderName != "" {
		fileName = fmt.Sprintf("%s%s/", fileName, folderName)
	}

	if fileName != "" {
		fileName = fmt.Sprintf("%s%s", fileName, filename)
	}

	return f.awsSVC.GetFile(fileName)
}

func (f fileService) ListFilesInFolder(username, folderName string) ([]string, *errorEntity.ServiceError) {
	fileName := fmt.Sprintf("%s/", username)

	if folderName != "" {
		fileName = fmt.Sprintf("%s%s/", fileName, folderName)
	}

	return f.awsSVC.ListFilesInFolder(fileName)
}

func (f fileService) DeleteFolder(username, folderName string) *errorEntity.ServiceError {
	fileName := fmt.Sprintf("%s/", username)

	if folderName != "" {
		fileName = fmt.Sprintf("%s%s/", fileName, folderName)
	}

	return f.awsSVC.DeleteFolder(fileName)
}
