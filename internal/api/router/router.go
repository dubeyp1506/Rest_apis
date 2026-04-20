package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {
	mainMux := http.NewServeMux()

	// Mount sub-routers to their respective prefixes using StripPrefix
	mainMux.Handle("/teachers/", http.StripPrefix("/teachers", teachersRouter()))
	mainMux.Handle("/students/", http.StripPrefix("/students", studentsRouter()))
	mainMux.Handle("/execs/", http.StripPrefix("/execs", execRouter()))

	return mainMux
}
