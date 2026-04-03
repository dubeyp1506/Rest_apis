package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strings"
)

func Teachers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetTeacherHanlder(w, r)
	case "POST":
		AddTeacherHandler(w, r)
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
func isValidSortOrder(order string) bool {
	return (order == "asc" || order == "desc")
}
func isValidFiled(filed string, params map[string]string) bool {
	_, ok := params[filed]
	return ok
}

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"subject":    "subject",
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
	}

	for param, dbFiled := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " and " + dbFiled + "=?"
			args = append(args, value)
		}
	}
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {

		var sortParts []string
		for _, sortParam := range sortParams {
			parts := strings.Split(sortParam, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortOrder(order) || !isValidFiled(field, params) {
				continue
			}
			sortParts = append(sortParts, field+" "+order)

		}
		if len(sortParts) > 0 {
			query += " order by " + strings.Join(sortParts, ", ")
		}
	}

	fmt.Println(query)
	return query, args
}

func GetTeacherHanlder(w http.ResponseWriter, r *http.Request) {

	db_name := os.Getenv("DB_NAME")
	db, err := sqlconnect.ConnectDb(db_name)
	if err != nil {
		http.Error(w, "Error Connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idstr := strings.TrimSuffix(path, "/")
	fmt.Println(idstr)

	var args []interface{}

	query := "Select id, first_name, last_name, email, class, subject from teachers where 1=1"
	query, args = addFilters(r, query, args)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Database Query Execution error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	teacherList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			http.Error(w, "Error during scaning database", http.StatusInternalServerError)
			return
		}
		teacherList = append(teacherList, teacher)
	}

	// var teacher models.Teacher
	// err = db.QueryRow(query, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		http.Error(w, "Teacher Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database Query error ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(teacherList)
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {

	db_name := os.Getenv("DB_NAME")
	db, err := sqlconnect.ConnectDb(db_name)
	if err != nil {
		http.Error(w, "Error Connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newTeachers []models.Teacher

	// Read body to allow multiple decoding attempts or just peek
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Try to decode as array
	err = json.Unmarshal(bodyBytes, &newTeachers)
	if err != nil {
		// If array fails, try decoding as a single object
		var singleTeacher models.Teacher
		err = json.Unmarshal(bodyBytes, &singleTeacher)
		if err != nil {
			http.Error(w, "Invalid Request Body: expected object or array", http.StatusBadRequest)
			return
		}
		newTeachers = append(newTeachers, singleTeacher)
	}

	stmt, err := db.Prepare("Insert Into teachers (first_name,last_name, email, class, subject) Values(?,?,?,?,?)")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in preparing the sql query: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
			return
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last inserted id", http.StatusInternalServerError)
		}
		newTeacher.ID = int(lastId)
		fmt.Println(lastId)
		addedTeachers[i] = newTeacher
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
