package middleware

import (
	"log"
	"net/http"
	"time"
)

// wrapper to extend http response writer to expose 
// the status codes 
type WrappedWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrappedWriter) WriteHeader(statuscode int) {
	w.ResponseWriter.WriteHeader(statuscode)
	w.StatusCode = statuscode
}

// logging middleware to track status codes, the url path, and response latency
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &WrappedWriter{
			ResponseWriter: w,
			StatusCode: http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.StatusCode, r.Method, r.URL.Path, time.Since(start))
	})
}