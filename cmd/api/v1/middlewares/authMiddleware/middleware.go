package authMiddleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tosinjs/cloud-backup/internal/entity/responseEntity"
	"tosinjs/cloud-backup/internal/service/authService"
)

type authMiddleware struct {
	authSVC authService.AuthService
}

func New(authSVC authService.AuthService) authMiddleware {
	return authMiddleware{
		authSVC: authSVC,
	}
}

func (a authMiddleware) CheckAdminStatus() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		status := c.GetString("user-status")

		if status != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseEntity.BuildErrorResponseObject(
				http.StatusUnauthorized, "Unauthorized", c.FullPath(),
			))
			return
		}
		c.Next()
	}
	return fn
}

func (a authMiddleware) VerifyJWT() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseEntity.BuildErrorResponseObject(
				http.StatusUnauthorized, "Unauthorized", c.FullPath(),
			))
			return
		}

		authArray := strings.Split(authHeader, " ")
		if len(authArray) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseEntity.BuildErrorResponseObject(
				http.StatusUnauthorized, "Unauthorized", c.FullPath(),
			))
			return
		}

		authPayload, isValid := a.authSVC.ValidateJWT(authArray[1])

		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseEntity.BuildErrorResponseObject(
				http.StatusUnauthorized, "Unauthorized", c.FullPath(),
			))
			return
		}

		c.Set("username", authPayload.Username)
		c.Set("user-status", authPayload.Status)
		c.Next()
	}
	return fn
}
