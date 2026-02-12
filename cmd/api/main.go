package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	mw "restapi/internal/api/middlewares"
	"time"
)

type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello From root"))
	log.Println("Hello From root")
}

func teachers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Hello From teachers using GET"))
	case "POST":
		fmt.Println("Body: ", r.Body)
		fmt.Println("Query: ", r.URL.Query())
		fmt.Println("name: ", r.URL.Query().Get("name"))
		fmt.Println("age: ", r.URL.Query().Get("age"))

		err := r.ParseForm()
		if err != nil {
			fmt.Println("Error parsing form: ", err)
		}
		fmt.Println("Form: ", r.Form)
		fmt.Println("name: ", r.FormValue("name"))
		fmt.Println("age: ", r.FormValue("age"))
		w.Write([]byte("Hello From teachers using POST"))
	case "PUT":
		w.Write([]byte("Hello From teachers using PUT"))
	case "DELETE":
		w.Write([]byte("Hello From teachers using DELETE"))
	case "PATCH":
		w.Write([]byte("Hello From teachers using PATCH"))
	default:
		w.Write([]byte("Hello From teachers using other methods"))
	}
}

func student(w http.ResponseWriter, r *http.Request) {
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

func exec(w http.ResponseWriter, r *http.Request) {
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
func main() {
	port := ":3001"

	rl := mw.NewRateLimiter(5, time.Minute)

	router := http.NewServeMux()

	router.HandleFunc("/", root)

	router.HandleFunc("/teachers/", teachers)

	router.HandleFunc("/students/", student)

	router.HandleFunc("/execs/", exec)

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	hppOptions := &mw.HppOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		WhiteList:                   []string{"first name", "last name", "email", "phone no"},
	}

	server := &http.Server{
		Addr:    port,
		Handler: mw.Hpp(*hppOptions)(rl.RateLimiterMiddleware(mw.Compression(mw.ResponseTime(mw.SecurityHandler(mw.Cors(router)))))),
		// Handler:   router,
		// Handler:   middlewares.Cors(router),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Listening on port", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error Starting the server :", err)
	}
}
