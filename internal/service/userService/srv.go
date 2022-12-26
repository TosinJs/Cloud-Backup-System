package userService

import (
	"fmt"
	"github.com/google/uuid"
	"tosinjs/cloud-backup/internal/entity/errorEntity"
	"tosinjs/cloud-backup/internal/entity/userEntity"
	"tosinjs/cloud-backup/internal/repository/userRepo"
	"tosinjs/cloud-backup/internal/service/cryptoService"
)

type userService struct {
	repo      userRepo.UserRepository
	cryptoSVC cryptoService.CryptoService
}

type UserService interface {
	CreateUser(userEntity.UserSignUpReq) *errorEntity.ServiceError
	LoginUser(userEntity.UserLoginReq) *errorEntity.ServiceError
}

func New(userRepo userRepo.UserRepository, cryptoSVC cryptoService.CryptoService) UserService {
	return userService{
		repo:      userRepo,
		cryptoSVC: cryptoSVC,
	}
}

func (u userService) CreateUser(req userEntity.UserSignUpReq) *errorEntity.ServiceError {
	req.UserId = uuid.New().String()
	password, srvErr := u.cryptoSVC.Hash(req.Password)

	if srvErr != nil {
		return srvErr
	}

	req.Password = password

	return u.repo.CreateUser(req)
}

func (u userService) LoginUser(req userEntity.UserLoginReq) *errorEntity.ServiceError {
	password, repoErr := u.repo.LoginUser(req)

	if repoErr != nil {
		return repoErr
	}
	if !u.cryptoSVC.Compare(req.Password, password) {
		return errorEntity.BadRequestError("Invalid Login Credentials", fmt.Errorf("invalid login credentials"))
	}
	return nil
}
