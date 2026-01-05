package middleware

import (
	"net/http"
	"time"

	"github.com/cobyabrahams/hungr/logger"
	"github.com/gofrs/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func RequestLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID
		requestID := uuid.Must(uuid.NewV4()).String()[:8]
		ctx := logger.WithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		// Wrap response writer to capture status
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Process request
		next(wrapped, r)

		// Log request
		duration := time.Since(start)
		logger.Info(ctx, "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.status,
			"duration_ms", duration.Milliseconds(),
		)
	}
}

func CORS(next http.HandlerFunc, methods string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", methods)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
