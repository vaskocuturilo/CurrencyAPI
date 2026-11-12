package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(context gin.Context, w http.ResponseWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		c.Next()

		duration := time.Since(start)

		slog.Info("HTTP Request",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", wrapped.status,
			"duration", duration,
		)
		c.Next()
	}
}
