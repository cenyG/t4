package handlers

import (
	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type downloadFileHandler struct {
	useCase  usecase.DownloadFileUseCase
	fileRepo repo.FileRepository
}

func NewDownloadFileHandler(useCase usecase.DownloadFileUseCase, repo repo.FileRepository) Handler {
	return &downloadFileHandler{useCase, repo}
}

func (d *downloadFileHandler) Handle(c *gin.Context) {
	fileId := c.Param("id")
	id, err := strconv.ParseInt(fileId, 10, 64)
	if err != nil {
		failResponse(c, 500, fmt.Sprintf("error while parsing id (must be int): %v", err))
		return
	}

	file, err := d.fileRepo.GetFile(c, id)
	if err != nil {
		failResponse(c, 500, fmt.Sprintf("fail to get file from store: %v", err))
		return
	}

	// Set header
	c.Header("Content-Disposition", "attachment; filename="+file.Name)
	c.Header("Content-Type", "application/octet-stream")

	// Проходим по всем серверам и запрашиваем части файла
	err = d.useCase.Download(c, file, c.Writer)
	if err != nil {
		failResponse(c, 500, fmt.Sprintf("fail while file loading: %v", err))
		return
	}

	c.Status(http.StatusOK)
}
