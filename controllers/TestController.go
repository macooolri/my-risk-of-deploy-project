package controllers

import (
	"net/http"
)

// TestWarning всегда возвращает 403 Forbidden для тестирования warning метрик
func TestWarning(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"error": "Forbidden - This endpoint always returns 403 for testing"}`))
}

// TestError всегда возвращает 500 Internal Server Error для тестирования error метрик
func TestError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "Internal Server Error - This endpoint always returns 500 for testing"}`))
}
