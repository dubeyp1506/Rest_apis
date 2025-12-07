package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	start  time.Time
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.status == 0 {
		rw.status = code
	}
	duration := time.Since(rw.start)
	rw.Header().Set("X-Response-Time", duration.String())
	rw.ResponseWriter.WriteHeader(code)
}

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Receiving Request in ResponseTime")

		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			status:         0,
			start:          time.Now(),
		}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(wrappedWriter.start)
		fmt.Printf("Request %s %s completed in %s with status %d\n", r.Method, r.URL.Path, duration, wrappedWriter.status)
	})
}
