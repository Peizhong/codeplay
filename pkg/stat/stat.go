package stat

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Number of HTTP requests",
	}, []string{"url", "method", "code"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP request duration",
	}, []string{"url", "method"})
)

func init() {
	prometheus.MustRegister(httpRequestTotal, httpRequestDuration)
}

func RecordRequest(url, method string, code int) {
	httpRequestTotal.WithLabelValues(url, method, strconv.Itoa(code)).Inc()
}

func RecordRequestDuration(url, method string, seconds float64) {
	httpRequestDuration.WithLabelValues(url, method).Observe(seconds)
}
