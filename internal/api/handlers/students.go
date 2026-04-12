package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"
	"strconv"
	"strings"
)

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {

	studentList, err := sqlconnect.GetStudentsHandlerFunc(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(studentList)
}

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, utils.ErrorHandler(err, "Invalid Id Type").Error(), http.StatusBadRequest)
		return
	}
	Student, err := sqlconnect.GetStudentHandlerFunc(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Student)
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {

	newStudents, err := sqlconnect.AddStudentHandlerFunc(r)
	if err != nil {
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(newStudents),
		Data:   newStudents,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, utils.ErrorHandler(err, "The id is not an integer").Error(), http.StatusBadRequest)
		return
	}

	updatedStudents, err := sqlconnect.UpdateStudentFunc(r, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student Not Found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudents)
}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	err := sqlconnect.PatchStudentsHandlerFunc(r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "One or more Student not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PatchStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, utils.ErrorHandler(err, "The id is not an integer").Error(), http.StatusBadRequest)
		return
	}

	existingStudent, err := sqlconnect.PatchStudentHandlerFunc(id, r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingStudent)
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, utils.ErrorHandler(err, "Invalid Student ID").Error(), http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteStudentHandlerFunc(id)
	if err != nil {
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, "Student Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Student deleted successfully"}
	json.NewEncoder(w).Encode(response)
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	err := sqlconnect.DeleteStudentsHandlerFunc(r)
	if err != nil {
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
