package http

import (
	handlers2 "T4_test_case/internal/restserver/http/handlers"
	"T4_test_case/internal/restserver/usecase"
	"github.com/gin-gonic/gin"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, provider usecase.Provider) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Handlers
	uploadHandler := handlers2.NewUploadFileHandler(provider.GetUploadFileUseCase())
	downloadHandler := handlers2.NewDownloadFileHandler(provider.GetDownloadFileUseCase())

	// Routers
	group := handler.Group("/file")

	group.POST("/", uploadHandler.Handle)
	group.GET("/:id", downloadHandler.Handle)
}
