package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func teachersRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /", handlers.AddTeacherHandler)
	mux.HandleFunc("PATCH /", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /", handlers.DeleteTeachersHandler)

	mux.HandleFunc("GET /{id}", handlers.GetTeacherHandler)
	mux.HandleFunc("PUT /{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /{id}", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /{id}", handlers.DeleteTeacherHandler)

	mux.HandleFunc("GET /{id}/students", handlers.GetStudentsOfTeachers)
	mux.HandleFunc("GET /{id}/studentsCount", handlers.GetStudentCountOfTeacher)

	return mux
}
