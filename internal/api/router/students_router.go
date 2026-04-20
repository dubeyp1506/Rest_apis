package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func studentsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /", handlers.AddStudentHandler)
	mux.HandleFunc("PATCH /", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /", handlers.DeleteStudentsHandler)

	mux.HandleFunc("GET /{id}", handlers.GetStudentHandler)
	mux.HandleFunc("PUT /{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("PATCH /{id}", handlers.PatchStudentHandler)
	mux.HandleFunc("DELETE /{id}", handlers.DeleteStudentHandler)

	return mux

}
