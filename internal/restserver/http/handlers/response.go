package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
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
	slog.Error(fmt.Sprintf("[failResponse] error: %s", msg))
	c.AbortWithStatusJSON(code, fail{msg})
}

func successResponse(c *gin.Context, id int64, filename string) {
	select {
	case <-c.Request.Context().Done():
		slog.Error(fmt.Sprintf("[successResponse] request error: %v", c.Request.Context().Err()))
	default:
		c.JSON(200, success{
			Status:   "OK",
			ID:       id,
			Filename: filename,
		})
	}
}
