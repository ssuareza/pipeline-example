package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Define prometheus metrics.
var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Histogram of latencies for HTTP requests.",
		// Buckets: prometheus.DefBuckets, // Default buckets (0.005s, 0.01s, 0.025s, etc.)
	},
	[]string{"path", "method", "status"},
)

// Init initializes and registers the metrics.
func Init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)
}

// Custom ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// InstrumentHandler is middleware that instruments HTTP requests.
func InstrumentHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		// Increment the counter with labels
		httpRequests.WithLabelValues(r.URL.Path, r.Method, fmt.Sprint(rw.status)).Inc()

		// Measure the latency
		duration := time.Since(start).Seconds()
		// log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.status, duration)
		httpDuration.WithLabelValues(r.URL.Path, r.Method, fmt.Sprint(rw.status)).Observe(duration)
	})
}

// Health returns a 200 OK status to indicate the service is healthy.
func Health(w http.ResponseWriter, req *http.Request) {
	status := map[string]string{
		"status":  "ok",
		"message": "service is healthy",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
