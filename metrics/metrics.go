package metrics

import (
	"net/http"
	"strconv"
	"time"

	"ai_chatbot_llama3.2_1B/types"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests by method, route and status code.",
	}, []string{"method", "path", "status"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds.",
		Buckets: []float64{0.05, 0.1, 0.5, 1, 2, 5, 10, 30, 60},
	}, []string{"method", "path"})

	chatsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chat_requests_in_flight",
		Help: "Chat requests currently being answered.",
	})

	timeToFirstToken = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "chat_time_to_first_token_seconds",
		Help:    "Delay between receiving a chat request and streaming the first token.",
		Buckets: []float64{0.1, 0.25, 0.5, 1, 2, 5, 10},
	})

	tokensStreamed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_tokens_streamed_total",
		Help: "Total chunks (roughly tokens) streamed to clients.",
	})
)

// metricsMiddleware records a count and a duration for every request.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// chi's wrapper captures the status code AND keeps http.Flusher working.
		// A hand-written wrapper would silently break our streaming.
		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		// Use the route pattern ("/chat"), not the raw URL, so the label
		// doesn't explode into thousands of unique values.
		path := chi.RouteContext(r.Context()).RoutePattern()
		if path == "" {
			path = "unmatched"
		}
		status := ww.Status()
		if status == 0 {
			status = http.StatusOK
		}

		httpRequests.WithLabelValues(r.Method, path, strconv.Itoa(status)).Inc()
		httpDuration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
	})
}

func InitMetrics() *types.Metrics {
	return &types.Metrics{
		ChatsInFlight:    chatsInFlight,
		TimeToFirstToken: timeToFirstToken,
		TokensStreamed:   tokensStreamed,
	}
}