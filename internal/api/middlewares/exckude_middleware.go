package middlewares

import (
	"net/http"
)

// ExcludeMiddleware takes a target middleware and a list of paths to exclude.
// If the current request path matches an excluded path exactly, it skips the middleware.
func ExcludeMiddleware(middleware func(http.Handler) http.Handler, excludedPaths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Pre-wrap the next handler with the target middleware
		handlerWithMiddleware := middleware(next)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Check if the current route is in the exclusion list
			for _, path := range excludedPaths {
				// Exact match (e.g. "/login" == "/login")
				if r.URL.Path == path {
					// Path is excluded! Skip the target middleware entirely and go straight to the next handler.
					next.ServeHTTP(w, r)
					return
				}
			}

			// 2. Path is NOT excluded. Run the target middleware.
			handlerWithMiddleware.ServeHTTP(w, r)
		})
	}
}
