package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("/", handlers.Root)

	router.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	router.HandleFunc("POST /teachers/", handlers.AddTeacherHandler)
	router.HandleFunc("PATCH /teachers/", handlers.PatchTeachersHandler)
	router.HandleFunc("DELETE /teachers/", handlers.DeleteTeachersHandler)

	router.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)
	router.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	router.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
	router.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	router.HandleFunc("/students/", handlers.Student)

	router.HandleFunc("/execs/", handlers.Exec)

	return router

}
