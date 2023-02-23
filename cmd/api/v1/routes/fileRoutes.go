package routes

import (
	"tosinjs/cloud-backup/cmd/api/v1/handlers/fileHandler"
	"tosinjs/cloud-backup/cmd/api/v1/middlewares/authMiddleware"
	"tosinjs/cloud-backup/cmd/api/v1/middlewares/limiterMiddleware"
	"tosinjs/cloud-backup/internal/service/authService"
	"tosinjs/cloud-backup/internal/service/fileService"

	"github.com/gin-gonic/gin"
)

func FileRoutes(
	v1 *gin.RouterGroup,
	fileSVC fileService.FileService,
	authSVC authService.AuthService,
) {
	authMiddyWare := authMiddleware.New(authSVC)
	limiterMiddyWare := limiterMiddleware.New()

	fileHandler := fileHandler.NewHandler(fileSVC)

	fileRoutes := v1.Group("/file")

	fileRoutes.Use(authMiddyWare.VerifyJWT())
	fileRoutes.Use(limiterMiddyWare.LimitRequests(10))

	fileRoutes.POST("",
		limiterMiddyWare.LimitRequests(1),
		limiterMiddyWare.LimitFileSize(200*1024*1024),
		fileHandler.UploadFile,
	)

	fileRoutes.GET("", fileHandler.GetFile)

	fileRoutes.DELETE("", fileHandler.DeleteFile)

	fileRoutes.GET("/list", fileHandler.ListFilesInFolder)

	fileRoutes.DELETE("/list", fileHandler.DeleteFolder)

	fileRoutes.PATCH("/flag", authMiddyWare.CheckAdminStatus(), fileHandler.FlagFile)
}
