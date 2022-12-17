package fileHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type fileHandler struct{}

func NewHandler() *fileHandler {
	return &fileHandler{}
}

func (f fileHandler) UploadFile(c *gin.Context) {
	c.JSON(http.StatusAccepted, "uploaded")
}

func (f fileHandler) GetFile(c *gin.Context) {
	c.JSON(http.StatusAccepted, "downloaded")
}

func (f fileHandler) DeleteFile(c *gin.Context) {
	c.JSON(http.StatusAccepted, "deleted")
}
