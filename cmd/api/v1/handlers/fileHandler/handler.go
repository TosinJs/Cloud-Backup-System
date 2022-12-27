package fileHandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"tosinjs/cloud-backup/internal/entity/responseEntity"
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
	username := c.GetString("username")

	folderName := c.Query("folder")

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responseEntity.BuildErrorResponseObject(
			http.StatusBadRequest, "No File Sent With This Request", c.FullPath(),
		))
		return
	}

	if srvErr := f.fileSVC.UploadFile(file, username, folderName); srvErr != nil {
		c.AbortWithStatusJSON(srvErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(srvErr, c.FullPath()))
		return
	}

	c.JSON(http.StatusCreated, responseEntity.BuildResponseObject(http.StatusCreated, c.FullPath(), nil))
}

func (f fileHandler) GetFile(c *gin.Context) {
	username := c.GetString("username")
	folderName := c.Query("folder")
	fileName := c.Query("filename")

	data, srvErr := f.fileSVC.GetFile(username, folderName, fileName)
	if srvErr != nil {
		c.AbortWithStatusJSON(srvErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(srvErr, c.FullPath()))
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "binary/octet-stream", data)
}

func (f fileHandler) DeleteFile(c *gin.Context) {
	username := c.GetString("username")
	folderName := c.Query("folder")
	fileName := c.Query("filename")

	if srvErr := f.fileSVC.DeleteFile(username, folderName, fileName); srvErr != nil {
		c.AbortWithStatusJSON(srvErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(srvErr, c.FullPath()))
		return
	}

	c.JSON(http.StatusAccepted, responseEntity.BuildResponseObject(http.StatusAccepted, c.FullPath(), nil))
}

func (f fileHandler) DeleteFolder(c *gin.Context) {
	username := c.GetString("username")
	folderName := c.Query("folder")

	if srvErr := f.fileSVC.DeleteFolder(username, folderName); srvErr != nil {
		c.AbortWithStatusJSON(srvErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(srvErr, c.FullPath()))
		return
	}

	c.JSON(http.StatusAccepted, responseEntity.BuildResponseObject(http.StatusAccepted, c.FullPath(), nil))
}

func (f fileHandler) ListFilesInFolder(c *gin.Context) {
	username := c.GetString("username")
	folderName := c.Query("folder")

	fileStructure, srvErr := f.fileSVC.ListFilesInFolder(username, folderName)
	if srvErr != nil {
		c.AbortWithStatusJSON(srvErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(srvErr, c.FullPath()))
		return
	}

	c.JSON(http.StatusOK, responseEntity.BuildResponseObject(http.StatusOK, c.FullPath(), fileStructure))
}
