package userHandler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tosinjs/cloud-backup/internal/entity/responseEntity"
	"tosinjs/cloud-backup/internal/entity/userEntity"
	"tosinjs/cloud-backup/internal/service/userService"
	"tosinjs/cloud-backup/internal/service/validationService"
)

type userHandler struct {
	userSVC       userService.UserService
	validationSVC validationService.ValidationService
}

func NewHandler(userSVC userService.UserService, validationSVC validationService.ValidationService) userHandler {
	return userHandler{
		userSVC:       userSVC,
		validationSVC: validationSVC,
	}
}

func (u userHandler) CreateUser(c *gin.Context) {
	var req userEntity.UserSignUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responseEntity.BuildErrorResponseObject(
			http.StatusBadRequest, err.Error(), c.FullPath(),
		))
		return
	}

	//validate input here
	if svcErr := u.validationSVC.Validate(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}

	if svcErr := u.userSVC.CreateUser(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}

	c.JSON(http.StatusCreated, responseEntity.BuildResponseObject(http.StatusCreated, c.FullPath(), nil))
}

func (u userHandler) LoginUser(c *gin.Context) {
	var req userEntity.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responseEntity.BuildErrorResponseObject(
			http.StatusBadRequest, err.Error(), c.FullPath(),
		))
		return
	}

	//validate input here
	if svcErr := u.validationSVC.Validate(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}

	if svcErr := u.userSVC.LoginUser(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}

	c.JSON(http.StatusCreated, responseEntity.BuildResponseObject(http.StatusCreated, c.FullPath(), nil))

}
