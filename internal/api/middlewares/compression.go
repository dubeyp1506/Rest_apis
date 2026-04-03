package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func Compression(next http.Handler) http.Handler {
	fmt.Println("This is a compression middleware ")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("===> Compression Middleware START")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			fmt.Println("<=== Compression Middleware END (no gzip)")
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// wrap the ResponseWriter

		w = &gzipResponseWriter{ResponseWriter: w, Writer: gz}

		next.ServeHTTP(w, r)
		fmt.Println("<=== Compression Middleware END")
	})
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
