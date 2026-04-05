package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strconv"
	"strings"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	teacherList, err := sqlconnect.GetTeachersHandlerFunc(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(teacherList)
}

func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Id Type", http.StatusBadRequest)
		return
	}
	teacher, err := sqlconnect.GetTeacherHandlerFunc(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(teacher)
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {

	newTeachers, err := sqlconnect.AddTeacherHandlerFunc(r)
	if err != nil {
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(newTeachers),
		Data:   newTeachers,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Printf("Error converting id to integer: %v\n", err)
		http.Error(w, "The id is not an integer", http.StatusBadRequest)
		return
	}

	updatedTeachers, err := sqlconnect.UpdateTeacherFunc(r, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher Not Found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeachers)
}

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	err := sqlconnect.PatchTeachersHandlerFunc(r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "One or more teacher not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Printf("Error converting id to integer: %v\n", err)
		http.Error(w, "The id is not an integer", http.StatusBadRequest)
		return
	}

	existingTeacher, err := sqlconnect.PatchTeacherHandlerFunc(id, r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteTeacherHandlerFunc(id)
	if err != nil {
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, "Teacher Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Teacher deleted successfully"}
	json.NewEncoder(w).Encode(response)
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	err := sqlconnect.DeleteTeachersHandlerFunc(r)
	if err != nil {
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
