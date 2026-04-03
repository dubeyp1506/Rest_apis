package utils

import "net/http"

type Middleware func(http.Handler) http.Handler

func ApplyMiddleWare(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
