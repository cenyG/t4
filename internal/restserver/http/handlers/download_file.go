package handlers

import (
	"fmt"
	"strconv"

	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type downloadFileHandler struct {
	useCase  usecase.DownloadFileUseCase
	fileRepo repo.FileRepository
}

func NewDownloadFileHandler(useCase usecase.DownloadFileUseCase, repo repo.FileRepository) Handler {
	return &downloadFileHandler{useCase, repo}
}

// Handle - download handler
func (d *downloadFileHandler) Handle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

	// Start file downloading
	err = d.useCase.Download(c, file, c.Writer)
	if err != nil {
		failResponse(c, 500, fmt.Sprintf("fail while file loading: %v", err))
		return
	}

	c.Status(http.StatusOK)
}
