package handlers

import "net/http"

func Student(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Hello From teachers using GET"))
	case "POST":
		w.Write([]byte("Hello From student using POST"))
	case "PUT":
		w.Write([]byte("Hello From student using PUT"))
	case "DELETE":
		w.Write([]byte("Hello From student using DELETE"))
	case "PATCH":
		w.Write([]byte("Hello From student using PATCH"))
	default:
		w.Write([]byte("Hello From student using other methods"))
	}
}
