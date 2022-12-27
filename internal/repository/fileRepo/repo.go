package fileRepo

import "tosinjs/cloud-backup/internal/entity/errorEntity"

type FileRepository interface {
	UploadFile(username, filename string) *errorEntity.ServiceError
	DeleteFile(filename string) *errorEntity.ServiceError
	FlagFile(filename string) (int, *errorEntity.ServiceError)
}
