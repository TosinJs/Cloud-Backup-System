package userRepo

import (
	"tosinjs/cloud-backup/internal/entity/errorEntity"
	"tosinjs/cloud-backup/internal/entity/userEntity"
)

type UserRepository interface {
	CreateUser(req userEntity.UserSignUpReq) *errorEntity.ServiceError
	LoginUser(req userEntity.UserLoginReq) (string, *errorEntity.ServiceError)
}
