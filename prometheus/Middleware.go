package prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаем wrapper для отслеживания status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Вызвать следующий обработчик
		next.ServeHTTP(rw, r)

		// После обработки запроса увеличить счетчик
		RequestCounter.With(prometheus.Labels{
			"path":   r.URL.Path,
			"method": r.Method,
			"status": fmt.Sprintf("%d", rw.statusCode),
		}).Inc()

		// Инкрементировать warning счетчик для 4xx ошибок
		if rw.statusCode >= 400 && rw.statusCode < 500 {
			WarningCounter.With(prometheus.Labels{
				"error_class": fmt.Sprintf("%dxx", rw.statusCode/100),
			}).Inc()
		}

		// Инкрементировать error счетчик для 5xx ошибок
		if rw.statusCode >= 500 && rw.statusCode < 600 {
			ErrorCounter.With(prometheus.Labels{
				"error_class": fmt.Sprintf("%dxx", rw.statusCode/100),
			}).Inc()
		}
	})
}
