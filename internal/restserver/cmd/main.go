package main

import (
	"T4_test_case/config"
	"T4_test_case/internal/restserver/http"
	"T4_test_case/internal/restserver/services"
	"T4_test_case/internal/restserver/usecase"
	"T4_test_case/pkg/httpserver"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Services
	storageServersProvider, err := services.NewStorageServersProvider(ctx)
	if err != nil {
		slog.Error("[main] services.NewStorageServersProvider() error: %v", err)
	}

	// UseCase
	useCaseProvider := usecase.NewUseCaseProvider(storageServersProvider)

	// HTTP Server
	handler := gin.New()
	http.NewRouter(handler, useCaseProvider)
	httpServer := httpserver.New(handler, httpserver.Port(config.Cfg.Rest.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Error("[main] signal %s", s)
	case err = <-httpServer.Notify():
		slog.Error("[main] httpServer.Notify: %v", err)
	}

	// Shutdown
	cancel()
	err = httpServer.Shutdown()
	if err != nil {
		slog.Error("[main] httpServer.Shutdown: %v", err)
	}
}
