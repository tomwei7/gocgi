package gocgi

import (
	"log"
	"net/http"
	"time"
)

type LoggerHandler struct {
	logger  *log.Logger
	handler http.Handler
}

type loggerResponseWriter struct {
	http.ResponseWriter

	statusCode int
	writeBytes int
}

func WithLogger(logger *log.Logger, handler http.Handler) http.Handler {
	return &LoggerHandler{logger: logger, handler: handler}
}

func (l *loggerResponseWriter) Write(b []byte) (int, error) {
	n, err := l.ResponseWriter.Write(b)
	l.writeBytes += n
	return n, err
}

func (l *loggerResponseWriter) WriteHeader(statusCode int) {
	l.statusCode = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func (h *LoggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lw := &loggerResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	h.handler.ServeHTTP(lw, r)
	duration := time.Now().Sub(start)
	h.logger.Printf("%s %s %s %s %d %d", r.RemoteAddr, r.Method, r.RequestURI, duration, lw.statusCode, lw.writeBytes)
}
