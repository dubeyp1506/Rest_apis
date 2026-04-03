package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("/", handlers.Root)

	router.HandleFunc("/teachers/", handlers.Teachers)

	router.HandleFunc("/students/", handlers.Student)

	router.HandleFunc("/execs/", handlers.Exec)

	return router

}
