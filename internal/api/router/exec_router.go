package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func execRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.GetExecsHandler)
	mux.HandleFunc("POST /", handlers.AddExecsHandler)
	mux.HandleFunc("PATCH /", handlers.PatchExecsHandler)

	mux.HandleFunc("GET /{id}", handlers.GetExecHandler)
	mux.HandleFunc("PUT /{id}", handlers.UpdateExecHandler)
	mux.HandleFunc("PATCH /{id}", handlers.PatchExecHandler)
	mux.HandleFunc("DELETE /{id}", handlers.DeleteExecHandler)
	mux.HandleFunc("POST /{id}/changePassword", handlers.UpdateExecHandler) // Placeholder

	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)                              // Placeholder
	mux.HandleFunc("POST /forgotPassword", handlers.UpdateExecHandler)                  // Placeholder
	mux.HandleFunc("POST /resetPassword/reset/{resetCode}", handlers.UpdateExecHandler) // Placeholder

	return mux

}
