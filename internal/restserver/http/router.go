package http

import (
	handlers2 "T4_test_case/internal/restserver/http/handlers"
	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/usecase"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"time"
)

// NewRouter .
func NewRouter(handler *gin.Engine, provider usecase.Provider, fileRepo repo.FileRepository) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          300 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	// Handlers
	uploadHandler := handlers2.NewUploadFileHandler(provider.GetUploadFileUseCase())
	downloadHandler := handlers2.NewDownloadFileHandler(provider.GetDownloadFileUseCase(), fileRepo)

	// Routers
	group := handler.Group("/files")

	group.POST("/", uploadHandler.Handle)
	group.GET("/:id", downloadHandler.Handle)
}
