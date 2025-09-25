package controllers

import (
	"go-app/db"
	"go-app/prometheus"
	u "go-app/utils"
	"net/http"
	"time"
)

// checkHealthStatus выполняет проверку здоровья и обновляет метрики
func checkHealthStatus() (map[string]interface{}, error) {
	database := db.GetDB()

	healthStatus := map[string]interface{}{
		"status":    "OK",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "notes-service",
		"version":   "1.0.0",
	}

	// Проверяем подключение к БД
	var err error
	if database != nil {
		err = database.DB().Ping()
	}

	if err != nil {
		healthStatus["status"] = "ERROR"
		healthStatus["database"] = "disconnected"
		healthStatus["error"] = err.Error()

		// Устанавливаем метрику службы в 0 (unhealthy)
		prometheus.HealthStatus.WithLabelValues("notes-service", "disconnected").Set(0)
		// Устанавливаем метрику БД в 0 (disconnected)
		prometheus.DatabaseStatus.Set(0)
	} else {
		healthStatus["database"] = "connected"

		// Устанавливаем метрику службы в 1 (healthy)
		prometheus.HealthStatus.WithLabelValues("notes-service", "connected").Set(1)
		// Устанавливаем метрику БД в 1 (connected)
		prometheus.DatabaseStatus.Set(1)
	}

	return healthStatus, err
}

// StartHealthMonitor запускает периодическую проверку здоровья
func StartHealthMonitor(interval time.Duration) {
	// Выполняем первую проверку сразу
	checkHealthStatus()

	// Создаем тикер для периодических проверок
	ticker := time.NewTicker(interval)

	// Бесконечный цикл проверок
	for range ticker.C {
		checkHealthStatus()
	}
}

var HealthCheck = func(w http.ResponseWriter, r *http.Request) {
	// Используем общую функцию проверки
	healthStatus, err := checkHealthStatus()

	// Устанавливаем правильный HTTP статус
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	u.Respond(w, healthStatus)
}
