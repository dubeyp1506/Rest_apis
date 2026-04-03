package handlers

import "net/http"

func Exec(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Hello From teachers using GET"))
	case "POST":
		w.Write([]byte("Hello From exec using POST"))
	case "PUT":
		w.Write([]byte("Hello From exec using PUT"))
	case "DELETE":
		w.Write([]byte("Hello From exec using DELETE"))
	case "PATCH":
		w.Write([]byte("Hello From exec using PATCH"))
	default:
		w.Write([]byte("Hello From exec using other methods"))
	}
}
