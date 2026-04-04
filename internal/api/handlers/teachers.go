package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strconv"
	"strings"
)

func Teachers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTeacherHanlder(w, r)
	case "POST":
		addTeacherHandler(w, r)
	case "PUT":
		updateTeacherHandler(w, r)
	case "DELETE":
		deleteTeacherHandler(w, r)
	case "PATCH":
		patchTeacherHandler(w, r)
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

func getTeacherHanlder(w http.ResponseWriter, r *http.Request) {

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

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {

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

func updateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	fmt.Println(id)
	if err != nil {
		fmt.Printf("Error converting id to integer: %v\n", err)
		http.Error(w, "The id is not an integer", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb(os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		http.Error(w, "Unable to Connect to Database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var updatedTeachers models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeachers)
	if err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	var existingTeacher models.Teacher

	err = db.QueryRow("Select id, first_name, last_name, email, subject, class from teachers where id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Subject,
		&existingTeacher.Class)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Teacher with id %d not found: %v\n", id, err)
			http.Error(w, "No Data found for this teacher", http.StatusNotFound)
			return
		}
		fmt.Printf("Error retrieving teacher data from DB: %v\n", err)
		http.Error(w, "unable to retrive Data", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("update teachers set first_name = ? , last_name = ? , email = ? , class = ? , subject = ? where id = ? ",
		updatedTeachers.FirstName, updatedTeachers.LastName, updatedTeachers.Email, updatedTeachers.Class, updatedTeachers.Subject, id)

	if err != nil {
		fmt.Printf("Error executing update query: %v\n", err)
		http.Error(w, "Error updating teachers", http.StatusInternalServerError)
		return
	}

	updatedTeachers.ID = id // Ensure ID is in the response
	fmt.Printf("Successfully updated teacher with ID %d\n", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeachers)
}

func patchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	fmt.Println(id)
	if err != nil {
		fmt.Printf("Error converting id to integer: %v\n", err)
		http.Error(w, "The id is not an integer", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb(os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		http.Error(w, "Unable to Connect to Database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("Select id, first_name, last_name, email, subject, class from teachers where id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Subject,
		&existingTeacher.Class)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Teacher with id %d not found: %v\n", id, err)
			http.Error(w, "No Data found for this teacher", http.StatusNotFound)
			return
		}
		fmt.Printf("Error retrieving teacher data from DB: %v\n", err)
		http.Error(w, "unable to retrive Data", http.StatusInternalServerError)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	val := reflect.ValueOf(&existingTeacher).Elem()
	typeOf := val.Type()
	fmt.Println("Type of :", typeOf)

	for k, v := range updates {
		for i := 0; i < val.NumField(); i++ {
			filed := typeOf.Field(i)
			tag := filed.Tag.Get("json")
			jsonKey := strings.Split(tag, ",")[0]

			if jsonKey == k {
				fieldVal := val.Field(i)
				if fieldVal.CanSet() {
					fieldVal.Set(reflect.ValueOf(v))
				}
			}

		}
	}

	_, err = db.Exec("update teachers set first_name = ? , last_name = ? , email = ? , class = ? , subject = ? where id = ? ",
		existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, id)

	if err != nil {
		fmt.Printf("Error executing update query: %v\n", err)
		http.Error(w, "Error updating teachers", http.StatusInternalServerError)
		return
	}

	existingTeacher.ID = id // Ensure ID is in the response
	fmt.Printf("Successfully updated teacher with ID %d\n", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}

func deleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb(os.Getenv("DB_NAME"))
	if err != nil {
		http.Error(w, "Failed to connect to Database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete teacher", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error retrieving deletion status", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Teacher deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
