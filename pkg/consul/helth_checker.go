package consul

import (
	"fmt"
	"net/http"
)

func RunHttpHealthCheck(port string) {
	http.HandleFunc("/health", healthHandler)
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

// Обработчик для маршрута /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
