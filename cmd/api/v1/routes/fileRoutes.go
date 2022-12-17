package routes

import (
	"github.com/gin-gonic/gin"
	"tosinjs/cloud-backup/cmd/api/v1/handlers/fileHandler"
)

func FileRoutes(v1 *gin.RouterGroup) {

	fileHandler := fileHandler.NewHandler()

	fileRoutes := v1.Group("/file")

	fileRoutes.POST("", fileHandler.UploadFile)

	fileRoutes.GET("", fileHandler.GetFile)

	fileRoutes.DELETE("", fileHandler.DeleteFile)
}
