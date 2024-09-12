package handlers

import (
	"T4_test_case/internal/restserver/usecase"
	"github.com/gin-gonic/gin"
)

type downloadFileHandler struct {
	useCase usecase.DownloadFileUseCase
}

func NewDownloadFileHandler(useCase usecase.DownloadFileUseCase) Handler {
	return &downloadFileHandler{useCase}
}

func (h *downloadFileHandler) Handle(c *gin.Context) {
	c.JSON(200, []int{1, 2, 3})
}
