package routes

import (
	"github.com/gin-gonic/gin"
	"tosinjs/cloud-backup/cmd/api/v1/handlers/fileHandler"
	"tosinjs/cloud-backup/cmd/api/v1/middlewares/authMiddleware"
	"tosinjs/cloud-backup/internal/service/authService"
	"tosinjs/cloud-backup/internal/service/fileService"
)

func FileRoutes(v1 *gin.RouterGroup, fileSVC fileService.FileService, authSVC authService.AuthService) {
	authMiddyWare := authMiddleware.New(authSVC)

	fileHandler := fileHandler.NewHandler(fileSVC)

	fileRoutes := v1.Group("/file")

	fileRoutes.Use(authMiddyWare.VerifyJWT())

	fileRoutes.POST("", fileHandler.UploadFile)

	fileRoutes.GET("", fileHandler.GetFile)

	fileRoutes.DELETE("", fileHandler.DeleteFile)

	fileRoutes.GET("/list", fileHandler.ListFilesInFolder)

	fileRoutes.DELETE("/list", fileHandler.DeleteFolder)

	fileRoutes.PATCH("/flag", authMiddyWare.CheckAdminStatus(), fileHandler.FlagFile)
}
