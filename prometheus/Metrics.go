package prometheus

import "github.com/prometheus/client_golang/prometheus"

var RequestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests received",
	},
	[]string{"path", "method", "status"},
)

var WarningCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_warnings_total",
		Help: "Total number of HTTP warnings, by class (4xx)",
	},
	[]string{"error_class"},
)

var ErrorCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP errors, by class (5xx)",
	},
	[]string{"error_class"},
)

// Business metrics
var NotesCreatedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "notes_created_total",
		Help: "Total number of notes created",
	},
)

var NotesUpdatedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "notes_updated_total",
		Help: "Total number of notes updated",
	},
)

var NotesDeletedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "notes_deleted_total",
		Help: "Total number of notes deleted",
	},
)

// Health status metric
var HealthStatus = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "app_health_status",
		Help: "Health status of the application (1 = healthy, 0 = unhealthy)",
	},
	[]string{"service", "database"},
)

// Database status metric
var DatabaseStatus = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "database_health_status",
		Help: "Database connection status (1 = connected, 0 = disconnected)",
	},
)
