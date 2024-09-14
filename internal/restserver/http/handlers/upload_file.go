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

// Handle - upload file handler
func (h *uploadFileHandler) Handle(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		failResponse(c, 500, err.Error())
		return
	}
	defer file.Close()

	// start uploading
	id, err := h.useCase.Upload(c, file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		failResponse(c, 500, err.Error())
		return
	}

	successResponse(c, id, fileHeader.Filename)
}
