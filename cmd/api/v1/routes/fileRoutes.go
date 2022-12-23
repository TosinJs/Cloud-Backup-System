package routes

import (
	"github.com/gin-gonic/gin"
	"tosinjs/cloud-backup/cmd/api/v1/handlers/fileHandler"
	"tosinjs/cloud-backup/internal/service/fileService"
)

func FileRoutes(v1 *gin.RouterGroup, fileSVC fileService.FileService) {

	fileHandler := fileHandler.NewHandler(fileSVC)

	fileRoutes := v1.Group("/file")

	fileRoutes.POST("", fileHandler.UploadFile)

	fileRoutes.GET("", fileHandler.GetFile)

	fileRoutes.DELETE("", fileHandler.DeleteFile)

	fileRoutes.GET("/list", fileHandler.ListFilesInFolder)

	fileRoutes.DELETE("/list", fileHandler.DeleteFolder)
}
