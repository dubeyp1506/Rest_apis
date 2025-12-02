package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
		path := r.URL.Path
		userId := strings.TrimPrefix(path, "/teachers/")
		fmt.Println(userId)
		fmt.Println(r.URL)
		fmt.Println(r)
		fmt.Println("Query Parameters", r.URL.Query())
		w.Write([]byte("Hello From teachers using GET"))
	case "POST":
		err := r.ParseForm() // means processing
		if err != nil {
			http.Error(w, "This is you err", http.StatusBadRequest)
			return
		}
		fmt.Println("row form data", r.Form)
		fmt.Println(r.PostForm)

		fmt.Println("-------")

		bodyRead, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "this is the error", http.StatusBadRequest)
			return
		}
		var user User
		err = json.Unmarshal(bodyRead, &user)

		fmt.Println(user)

		fmt.Println(string(bodyRead))
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
		err := r.ParseForm() // means processing
		if err != nil {
			http.Error(w, "This is you err", http.StatusBadRequest)
			return
		}
		fmt.Println("row form data", r.Form)
		fmt.Println(r.PostForm)

		fmt.Println("-------")

		bodyRead, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "this is the error", http.StatusBadRequest)
			return
		}
		var user User
		err = json.Unmarshal(bodyRead, &user)

		fmt.Println(user)

		fmt.Println(string(bodyRead))
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
		err := r.ParseForm() // means processing
		if err != nil {
			http.Error(w, "This is you err", http.StatusBadRequest)
			return
		}
		fmt.Println("row form data", r.Form)
		fmt.Println(r.PostForm)

		fmt.Println("-------")

		bodyRead, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "this is the error", http.StatusBadRequest)
			return
		}
		var user User
		err = json.Unmarshal(bodyRead, &user)

		fmt.Println(user)

		fmt.Println(string(bodyRead))
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
	port := ":3000"

	http.HandleFunc("/", root)

	http.HandleFunc("/teachers/", teachers)

	http.HandleFunc("/students/", student)

	http.HandleFunc("/execs/", exec)

	fmt.Println("Listening on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error Starting the server :", err)
	}
}
