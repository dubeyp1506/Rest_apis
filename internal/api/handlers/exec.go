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
	"time"
)

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	execsList, err := sqlconnect.GetExecsHandlerFunc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(execsList)
}

func GetExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "client: Invalid Id Type", http.StatusBadRequest)
		return
	}
	exec, err := sqlconnect.GetExecHandlerFunc(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "executive not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(exec)
}

func AddExecsHandler(w http.ResponseWriter, r *http.Request) {
	newExecs, err := sqlconnect.AddExecHandlerFunc(r)
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
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(newExecs),
		Data:   newExecs,
	}

	json.NewEncoder(w).Encode(response)
}

func AddExecHandler(w http.ResponseWriter, r *http.Request) {

	newExecs, err := sqlconnect.AddExecHandlerFunc(r)
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
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(newExecs),
		Data:   newExecs,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, utils.ErrorHandler(err, "The id is not an integer").Error(), http.StatusBadRequest)
		return
	}

	updatedExecs, err := sqlconnect.UpdateExecFunc(r, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Exec Not Found", http.StatusNotFound)
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
	json.NewEncoder(w).Encode(updatedExecs)
}

func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {
	err := sqlconnect.PatchExecsHandlerFunc(r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "One or more Exec not found", http.StatusNotFound)
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

func PatchExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, utils.ErrorHandler(err, "The id is not an integer").Error(), http.StatusBadRequest)
		return
	}

	existingExec, err := sqlconnect.PatchExecHandlerFunc(id, r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Exec not found", http.StatusNotFound)
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
	json.NewEncoder(w).Encode(existingExec)
}

func DeleteExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, utils.ErrorHandler(err, "Invalid Exec ID").Error(), http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteExecHandlerFunc(id)
	if err != nil {
		if strings.Contains(err.Error(), "client:") {
			http.Error(w, "Exec Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Exec deleted successfully"}
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Exec

	//Data Validation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	//Search for user into db
	db := sqlconnect.DB
	user := models.Exec{}
	err := db.QueryRow("Select id, first_name,last_name,email,username,password,is_active, role from execs where username=?", req.Username).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.IsActive, &user.Role,
	)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusBadRequest)
		return
	}

	// Is user active
	if !user.IsActive {
		http.Error(w, "user is not active", http.StatusForbidden)
		return
	}

	// Validate the Password

	isValidPassword := utils.VerifyPassword(user, req)
	if !isValidPassword {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//Generate Jwt Token

	token, err := utils.SignToken(user.ID, req.Username, user.Role)
	if err != nil {
		http.Error(w, "could not create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Content-type", "application/json")
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	json.NewEncoder(w).Encode(response)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                 // Good!
		SameSite: http.SameSiteLaxMode, // <--- ADD THIS
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out Successfully"}`))
}
