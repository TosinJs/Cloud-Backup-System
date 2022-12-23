package fileHandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"tosinjs/cloud-backup/internal/service/fileService"
)

type fileHandler struct {
	fileSVC fileService.FileService
}

func NewHandler(fileSVC fileService.FileService) *fileHandler {
	return &fileHandler{
		fileSVC: fileSVC,
	}
}

func (f fileHandler) UploadFile(c *gin.Context) {

	folderName := c.Query("folder")

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("No File Sent With This Request")
		return
	}

	if err = f.fileSVC.UploadFile(file, folderName); err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusAccepted, "uploaded")
}

func (f fileHandler) GetFile(c *gin.Context) {
	folderName := c.Query("folder")
	fileName := c.Query("filename")

	data, err := f.fileSVC.GetFile(folderName, fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "binary/octet-stream", data)
}

func (f fileHandler) DeleteFile(c *gin.Context) {
	folderName := c.Query("folder")
	fileName := c.Query("filename")

	if err := f.fileSVC.DeleteFile(folderName, fileName); err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusAccepted, "deleted")
}

func (f fileHandler) DeleteFolder(c *gin.Context) {
	folderName := c.Query("folder")

	err := f.fileSVC.DeleteFolder(folderName)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusAccepted, "deleted")
}

func (f fileHandler) ListFilesInFolder(c *gin.Context) {
	folderName := c.Query("folder")

	fileStructure, err := f.fileSVC.ListFilesInFolder(folderName)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, fileStructure)
}
