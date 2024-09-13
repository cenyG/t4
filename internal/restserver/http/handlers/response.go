package handlers

import (
	"github.com/gin-gonic/gin"
)

type fail struct {
	Error string `json:"error" example:"message"`
}

type success struct {
	Status   string `json:"status"`
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
}

func failResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, fail{msg})
}

func successResponse(c *gin.Context, id int64, filename string) {
	c.JSON(200, success{
		Status:   "OK",
		ID:       id,
		Filename: filename,
	})
}
