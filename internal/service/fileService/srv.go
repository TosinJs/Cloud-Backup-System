package fileService

import (
	"fmt"
	"mime/multipart"

	"tosinjs/cloud-backup/internal/entity/errorEntity"
	"tosinjs/cloud-backup/internal/repository/fileRepo"
	"tosinjs/cloud-backup/internal/service/awsService"
)

type fileService struct {
	awsSVC awsService.AWSService
	repo   fileRepo.FileRepository
}

type FileService interface {
	UploadFile(hFile *multipart.FileHeader, username, folderName string) *errorEntity.ServiceError
	DeleteFile(username, folderName, filename string) *errorEntity.ServiceError
	GetFile(username, folderName, filename string) ([]byte, *errorEntity.ServiceError)
	DeleteFolder(username, folderName string) *errorEntity.ServiceError
	FlagFile(filename string) *errorEntity.ServiceError
	ListFilesInFolder(username, folderName string) ([]string, *errorEntity.ServiceError)
}

func New(awsSVC awsService.AWSService, repo fileRepo.FileRepository) FileService {
	return fileService{
		awsSVC: awsSVC,
		repo:   repo,
	}
}

func (f fileService) FlagFile(filename string) *errorEntity.ServiceError {
	flagCount, svcErr := f.repo.FlagFile(filename)
	if svcErr != nil {
		return svcErr
	}
	if flagCount > 2 {
		svcErr = f.repo.DeleteFile(filename)
		if svcErr != nil {
			return svcErr
		}
		svcErr = f.awsSVC.DeleteFile(filename)
		if svcErr != nil {
			return svcErr
		}
	}
	return nil
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

	if svcErr := f.repo.UploadFile(username, filename); svcErr != nil {
		return svcErr
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

	if svcErr := f.repo.DeleteFile(fileName); svcErr != nil {
		return svcErr
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
