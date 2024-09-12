package handlers

import (
	"T4_test_case/internal/restserver/usecase"
	"github.com/gin-gonic/gin"
)

type uploadFileHandler struct {
	useCase usecase.UploadFileUseCase
}

func NewUploadFileHandler(useCase usecase.UploadFileUseCase) Handler {
	return &uploadFileHandler{useCase}
}

func (h *uploadFileHandler) Handle(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("fileHeader")
	if err != nil {
		c.Error(err)
		return
	}
	defer file.Close()

	err = h.useCase.Upload(c, file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		c.Error(err)
	}
}
