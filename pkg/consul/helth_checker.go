package consul

import (
	"fmt"
	"log/slog"
	"net/http"
)

func RunHttpHealthCheck(port string) {
	http.HandleFunc("/health", healthHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		slog.Error(fmt.Sprintf("health checker server error: %v", err))
	}
}

// Обработчик для маршрута /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		slog.Error(fmt.Sprintf("healthHandler error: %v", err))
	}
}
