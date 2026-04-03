package middlewares

import (
	"fmt"
	"net/http"
)

type HppOptions struct {
	CheckQuery                  bool
	CheckBody                   bool
	CheckBodyOnlyForContentType string
	WhiteList                   []string
}

func Hpp(options HppOptions) func(http.Handler) http.Handler {
	fmt.Println("This is a Hpp middleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("===> HPP Middleware START")

			if options.CheckBody && r.Method == http.MethodPost && isValid(options.CheckBodyOnlyForContentType, r.Header.Get("Content-Type")) {
				filterBodyParamsThroughWhiteList(r, options.WhiteList)
			}

			if options.CheckQuery && r.URL.Query() != nil {
				filterQueryParamsThroughWhiteList(r, options.WhiteList)
			}

			next.ServeHTTP(w, r)
			fmt.Println("<=== HPP Middleware END")
		})
	}
}

func isValid(contentType string, headerValue string) bool {
	return contentType == headerValue
}

func filterBodyParamsThroughWhiteList(r *http.Request, whiteList []string) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range r.Form {
		if len(v) > 1 {
			r.Form.Set(k, v[0])
		}
		if !isInWhiteList(whiteList, k) {
			delete(r.Form, k)
		}
	}
}

func filterQueryParamsThroughWhiteList(r *http.Request, whiteList []string) {
	query := r.URL.Query()

	for k, v := range query {
		if len(v) > 1 {
			query.Set(k, v[0])
		}
		if !isInWhiteList(whiteList, k) {
			query.Del(k)
		}
	}
	r.URL.RawQuery = query.Encode()
}

func isInWhiteList(whiteList []string, key string) bool {
	for _, v := range whiteList {
		if v == key {
			return true
		}
	}
	return false
}
