package logging

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	Environment string
	Output      io.Writer
}

// NewLogger creates a configured zerolog logger
func NewLogger(output io.Writer, env string) zerolog.Logger {
	// Set global log level based on environment
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Configure logger with pretty printing in development
	if env == "development" && output == os.Stdout {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
	}

	return zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()
}

// LogRequest is standard http.Handler middleware that logs HTTP requests using zerolog
func LogRequest(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200} // default status 200
			next.ServeHTTP(lrw, r)

			var logEvent *zerolog.Event
			if lrw.statusCode >= 400 {
				logEvent = logger.Error()
			} else {
				logEvent = logger.Info()
			}

			logEvent.
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", lrw.statusCode).
				Str("ip", r.RemoteAddr).
				Str("user-agent", r.UserAgent()).
				Dur("duration", time.Since(start)).
				Msg("request processed")
		})
	}
}

// LoggingMiddleware returns a Gin middleware that logs HTTP requests using zerolog
func LoggingMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Wrap gin.ResponseWriter to capture status code
		lrw := &ginLoggingResponseWriter{ResponseWriter: c.Writer, statusCode: 200}
		c.Writer = lrw

		// Process request
		c.Next()

		duration := time.Since(start)

		event := logger.Info()
		if lrw.statusCode >= 400 {
			event = logger.Error()
		}

		event.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", lrw.statusCode).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Dur("duration", duration).
			Msg("request processed")
	}
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code for standard http.Handler
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code for standard http.Handler
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// ginLoggingResponseWriter wraps gin.ResponseWriter to capture status code for Gin middleware
type ginLoggingResponseWriter struct {
	gin.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code for Gin middleware
func (lrw *ginLoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Initialize global logger on package import
func init() {
	log.Logger = NewLogger(os.Stdout, "development")
}
