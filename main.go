package main

import (
	"fmt"
	"go-app/controllers"
	"go-app/prometheus"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var LogPath = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: %s (%s)", r.Host, r.RequestURI, r.Method))
		next.ServeHTTP(w, r)
	})
}

func init() {
	prom.MustRegister(prometheus.RequestCounter)
	prom.MustRegister(prometheus.WarningCounter)
	prom.MustRegister(prometheus.ErrorCounter)

	// Регистрация бизнес-метрик
	prom.MustRegister(prometheus.NotesCreatedTotal)
	prom.MustRegister(prometheus.NotesUpdatedTotal)
	prom.MustRegister(prometheus.NotesDeletedTotal)

	// Регистрация метрики health status
	prom.MustRegister(prometheus.HealthStatus)
	prom.MustRegister(prometheus.DatabaseStatus)
}

func main() {
	router := mux.NewRouter()

	// Запускаем периодический мониторинг здоровья (каждые 10 секунд)
	go controllers.StartHealthMonitor(10 * time.Second)
	log.Info("Started health monitor with 10s interval")

	// Применяем middleware перед определением маршрутов
	router.Use(LogPath)
	router.Use(prometheus.MetricsMiddleware)

	router.HandleFunc("/notes", controllers.NoteQuery).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/notes", controllers.NoteCreate).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteRetrieve).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteUpdate).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteDelete).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/health", controllers.HealthCheck).Methods(http.MethodGet, http.MethodOptions)

	// Test endpoints для проверки метрик
	router.HandleFunc("/warning", controllers.TestWarning).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/error", controllers.TestError).Methods(http.MethodGet, http.MethodOptions)

	router.Handle("/metrics", promhttp.Handler())

	log.Info("Listening on 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
