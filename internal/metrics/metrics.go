package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	HttpRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	HttpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "route", "status"},
	)

	ActiveWebsocketConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_websocket_connections",
			Help: "Current number of active connections",
		},
	)

	DbRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_requests_total",
			Help: "Total number of DB queries",
		},
		[]string{"query"},
	)

	DbQueryDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
		},
		[]string{"query"},
	)

	DbErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total number of DB errors",
		},
		[]string{"query"},
	)
)

func Init() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDurationSeconds,
		HttpErrorsTotal,
		ActiveWebsocketConnections,
		DbRequestsTotal,
		DbQueryDurationSeconds,
		DbErrorsTotal,
	)
}
