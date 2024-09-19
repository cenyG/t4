package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"T4_test_case/config"
	"T4_test_case/internal/restserver/http"
	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/services"
	"T4_test_case/internal/restserver/usecase"
	"T4_test_case/pkg/db"
	"T4_test_case/pkg/httpserver"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	ctx, cancel := context.WithCancel(context.Background())

	// Services
	storageServersProvider, err := services.NewStorageServersProvider(ctx)
	if err != nil {
		log.Fatalf("[main] services.NewStorageServersProvider() error: %v", err)
	}

	// Store
	store, err := db.NewStore()
	if err != nil {
		log.Fatalf("[main] db.NewStore() error: %v", err)
	}
	fileRepo := repo.NewFileRepository(store)

	// UseCase
	useCaseProvider := usecase.NewUseCaseProvider(storageServersProvider, fileRepo)

	// HTTP Server
	handler := gin.New()
	http.NewRouter(handler, useCaseProvider, fileRepo)
	httpServer := httpserver.New(handler, httpserver.Port(config.Cfg.Rest.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Error(fmt.Sprintf("[main] signal %s", s))
		cancel()
	case err = <-httpServer.Notify():
		slog.Error(fmt.Sprintf("[main] httpServer.Notify: %v", err))
		cancel()
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		slog.Error(fmt.Sprintf("[main] httpServer.Shutdown: %v", err))
	}
}
