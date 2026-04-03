package handlers

import (
	"log"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello From root"))
	log.Println("Hello From root")
}
