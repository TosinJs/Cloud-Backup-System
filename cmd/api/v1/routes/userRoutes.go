package routes

import (
	"github.com/gin-gonic/gin"
	"tosinjs/cloud-backup/cmd/api/v1/handlers/userHandler"
	"tosinjs/cloud-backup/internal/service/userService"
	"tosinjs/cloud-backup/internal/service/validationService"
)

func UserRoutes(
	v1 *gin.RouterGroup,
	userSVC userService.UserService,
	validationSVC validationService.ValidationService,
) {

	userHandler := userHandler.NewHandler(userSVC, validationSVC)

	userRoutes := v1.Group("/user")

	userRoutes.POST("/signup", userHandler.CreateUser)

	userRoutes.POST("/login", userHandler.LoginUser)
}

/*
user can login and signup
single admin is seeded from the backend
single admin can make others admin

authToken and Refresh Tokens
Credentials will be stored in the tokens

any user can create a post
each post id saved to a folder in an s3 bucket with the users name

we have two tables
one users table - stores userid, username, password, email
one media table - stores userid, username, folderpath, flagCount

a flagCount above three will trigger a delete operation on the user
*/
