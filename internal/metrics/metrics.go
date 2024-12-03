package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
		},
		[]string{"pattern", "method", "status"},
	)

	HttpRequestsCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_inflight_current",
		},
		[]string{},
	)

	HttpRequestsInflightMax = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_inflight_max",
		},
		[]string{},
	)

	HttpRequestsDurationHistorgram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds_historgram",
			Buckets: []float64{
				0.1,  // 100 ms
				0.2,  // 200 ms
				0.25, // 250 ms
				0.3,  // 300 ms
				0.5,  // 500 ms
				1,    // 1 s
				2,    // 2 s
				3,    // 3 s
				5,    // 5 s
				7,    // 7 s
			},
		},
		[]string{"pattern", "method"},
	)

	HttpRequestsDurationSummary = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_seconds_summary",
			Objectives: map[float64]float64{
				0.99: 0.001, // 0.99 +- 0.001
				0.95: 0.01,  // 0.95 +- 0.01
				0.5:  0.05,  // 0.5 +- 0.05
			},
		},
		[]string{"pattern", "method"},
	)

	DatabaseQuesriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
		},
		[]string{"query_type", "source"},
	)

	DbQueryErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_errors_total",
			Help: "Total number of repo query errors",
		},
		[]string{"query_type", "source", "error_type"},
	)

	DbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "db_query_duration_seconds",
			Help: "Histogram of repo query durations",
			Buckets: []float64{
				0.1,  // 100 ms
				0.25, // 250 ms
				0.3,  // 300 ms
				0.5,  // 500 ms
				1,    // 1 s
				2.5,  // 2 s 500 ms
				5,    // 5 s
				10,   // 10 s
			},
		},
		[]string{"query_type", "source"},
	)
)
