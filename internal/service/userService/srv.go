package userService

import (
	"fmt"
	"github.com/google/uuid"
	"tosinjs/cloud-backup/internal/entity/errorEntity"
	"tosinjs/cloud-backup/internal/entity/tokenEntity"
	"tosinjs/cloud-backup/internal/entity/userEntity"
	"tosinjs/cloud-backup/internal/repository/userRepo"
	"tosinjs/cloud-backup/internal/service/authService"
	"tosinjs/cloud-backup/internal/service/cryptoService"
)

type userService struct {
	repo      userRepo.UserRepository
	cryptoSVC cryptoService.CryptoService
	authSVC   authService.AuthService
}

type UserService interface {
	CreateUser(userEntity.UserSignUpReq) (*userEntity.UserSignUpRes, *errorEntity.ServiceError)
	LoginUser(userEntity.UserLoginReq) (*userEntity.UserLoginRes, *errorEntity.ServiceError)
}

func New(userRepo userRepo.UserRepository, cryptoSVC cryptoService.CryptoService, authSVC authService.AuthService) UserService {
	return userService{
		repo:      userRepo,
		cryptoSVC: cryptoSVC,
		authSVC:   authSVC,
	}
}

func (u userService) CreateUser(req userEntity.UserSignUpReq) (*userEntity.UserSignUpRes, *errorEntity.ServiceError) {
	req.UserId = uuid.New().String()
	password, srvErr := u.cryptoSVC.Hash(req.Password)

	if srvErr != nil {
		return nil, srvErr
	}

	req.Password = password

	srvErr = u.repo.CreateUser(req)
	if srvErr != nil {
		return nil, srvErr
	}

	token, srvErr := u.authSVC.GenerateJWT(60, tokenEntity.JWTPayload{
		Username: req.Username,
		Status:   "user",
	})
	if srvErr != nil {
		return nil, srvErr
	}

	return &userEntity.UserSignUpRes{Token: token}, nil
}

func (u userService) LoginUser(req userEntity.UserLoginReq) (*userEntity.UserLoginRes, *errorEntity.ServiceError) {
	userDetails, repoErr := u.repo.LoginUser(req)

	if repoErr != nil {
		return nil, repoErr
	}
	if !u.cryptoSVC.Compare(req.Password, userDetails.Password) {
		return nil, errorEntity.BadRequestError("Invalid Login Credentials", fmt.Errorf("invalid login credentials"))
	}
	token, srvErr := u.authSVC.GenerateJWT(60, tokenEntity.JWTPayload{
		Username: req.Username,
		Status:   userDetails.Status,
	})
	if srvErr != nil {
		return nil, srvErr
	}

	return &userEntity.UserLoginRes{Token: token}, nil
}
