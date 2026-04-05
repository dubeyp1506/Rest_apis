package sqlconnect

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"restapi/internal/models"
	"strconv"
	"strings"
)

func GetTeacherHandlerFunc(id int) (models.Teacher, error) {
	db_name := os.Getenv("DB_NAME")
	db, err := ConnectDb(db_name)
	if err != nil {
		return models.Teacher{}, err
	}
	defer db.Close()

	query := "Select id, first_name, last_name, email, class, subject from teachers where id=? "

	var teacher models.Teacher
	err = db.QueryRow(query, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		return models.Teacher{}, err
	} else if err != nil {
		return models.Teacher{}, err
	}
	return teacher, nil
}

func GetTeachersHandlerFunc(r *http.Request) ([]models.Teacher, error) {
	db_name := os.Getenv("DB_NAME")
	db, err := ConnectDb(db_name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var args []interface{}

	query := "Select id, first_name, last_name, email, class, subject from teachers where 1=1"
	query, args = AddFilters(r, query, args)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teacherList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return teacherList, err
		}
		teacherList = append(teacherList, teacher)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teacherList, nil
}

func AddTeacherHandlerFunc(r *http.Request) ([]models.Teacher, error) {
	db_name := os.Getenv("DB_NAME")
	db, err := ConnectDb(db_name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var newTeachers []models.Teacher

	// Read body to allow multiple decoding attempts or just peek
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Try to decode as array
	err = json.Unmarshal(bodyBytes, &newTeachers)
	if err != nil {
		// If array fails, try decoding as a single object
		var singleTeacher models.Teacher
		err = json.Unmarshal(bodyBytes, &singleTeacher)
		if err != nil {
			return nil, fmt.Errorf("client: Invalid Request Body: expected object or array")
		}
		newTeachers = append(newTeachers, singleTeacher)
	}

	stmt, err := db.Prepare("Insert Into teachers (first_name,last_name, email, class, subject) Values(?,?,?,?,?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			return nil, err
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		newTeacher.ID = int(lastId)
		fmt.Println(lastId)
		addedTeachers[i] = newTeacher
	}
	return newTeachers, nil
}

func DeleteTeachersHandlerFunc(r *http.Request) error {
	db, err := ConnectDb(os.Getenv("DB_NAME"))
	if err != nil {
		return err
	}
	defer db.Close()

	var deletes []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&deletes); err != nil {
		return err
	}
	var deleteIdList []interface{}

	for _, v := range deletes {
		if id, ok := v["id"].(float64); ok {
			deleteIdList = append(deleteIdList, int(id))
		} else if id, ok := v["id"].(int); ok {
			deleteIdList = append(deleteIdList, id)
		} else if idStr, ok := v["id"].(string); ok {
			if id, err := strconv.Atoi(idStr); err == nil {
				deleteIdList = append(deleteIdList, id)
			}
		}
	}

	if len(deleteIdList) == 0 {
		return fmt.Errorf("no valid IDs provided")
	}
	placeholders := make([]string, len(deleteIdList))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf("DELETE FROM teachers WHERE id IN (%s)", strings.Join(placeholders, ","))
	// 3. Use db.Exec with the list of IDs
	_, err = db.Exec(query, deleteIdList...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTeacherHandlerFunc(id int) error {
	db, err := ConnectDb(os.Getenv("DB_NAME"))
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("client: Teacher not found")
	}
	return nil
}

func PatchTeacherHandlerFunc(id int, r *http.Request) (models.Teacher, error) {
	db, err := ConnectDb(os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return models.Teacher{}, err
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
			return models.Teacher{}, err
		}
		fmt.Printf("Error retrieving teacher data from DB: %v\n", err)
		return models.Teacher{}, err
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		return models.Teacher{}, fmt.Errorf("client: Invalid Request Body")
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
					currVal := reflect.ValueOf(v)
					// Add this check like you did in the other function!
					if currVal.Type().ConvertibleTo(fieldVal.Type()) {
						fieldVal.Set(currVal.Convert(fieldVal.Type()))
					} else {
						return models.Teacher{}, fmt.Errorf("client: type mismatch for field %s", k)
					}
				}
			}

		}
	}

	_, err = db.Exec("update teachers set first_name = ? , last_name = ? , email = ? , class = ? , subject = ? where id = ? ",
		existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, id)

	if err != nil {
		fmt.Printf("Error executing update query: %v\n", err)
		return models.Teacher{}, err
	}

	existingTeacher.ID = id // Ensure ID is in the response
	fmt.Printf("Successfully updated teacher with ID %d\n", id)
	return existingTeacher, nil
}

func PatchTeachersHandlerFunc(r *http.Request) error {
	db, err := ConnectDb(os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return err
	}
	defer db.Close()

	var updates []map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		return fmt.Errorf("client: invalid JSON body: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, update := range updates {
		// Robust ID extraction: handles Numbers (float64), Integers, and Strings from JSON
		var id int
		var idValid bool
		if v, ok := update["id"].(float64); ok {
			id = int(v)
			idValid = true
		} else if v, ok := update["id"].(int); ok {
			id = v
			idValid = true
		} else if v, ok := update["id"].(string); ok {
			if parsedId, err := strconv.Atoi(v); err == nil {
				id = parsedId
				idValid = true
			}
		}

		if !idValid {
			tx.Rollback()
			return fmt.Errorf("client: Invalid or missing teacher id in update")
		}

		var teacherFromDb models.Teacher
		// Use tx instead of db to keep it in the transaction
		err = tx.QueryRow("Select id, first_name, last_name, email, class, subject from teachers where id = ?", id).Scan(
			&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class,
			&teacherFromDb.Subject,
		)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return err
			}
			return err
		}

		val := reflect.ValueOf(&teacherFromDb).Elem()
		typeOf := val.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < val.NumField(); i++ {
				filed := typeOf.Field(i)
				tag := filed.Tag.Get("json")
				jsonKey := strings.Split(tag, ",")[0]

				if jsonKey == k {
					fieldVal := val.Field(i)
					if fieldVal.CanSet() {
						currVal := reflect.ValueOf(v)
						// Fix: check if the VALUE (currVal) is convertible to the FIELD
						if currVal.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(currVal.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							return fmt.Errorf("client: Type mismatch for field %s", k)
						}
					}
				}

			}
		}
		_, err = tx.Exec("Update teachers set first_name = ?, last_name = ?, email = ?, class = ?, subject = ? where id = ?",
			teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, teacherFromDb.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateTeacherFunc(r *http.Request, id int) (models.Teacher, error) {
	db, err := ConnectDb(os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return models.Teacher{}, err
	}
	defer db.Close()

	var updatedTeachers models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeachers)
	if err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		return models.Teacher{}, fmt.Errorf("client: Invalid Request Body")
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
			return models.Teacher{}, err
		}
		fmt.Printf("Error retrieving teacher data from DB: %v\n", err)
		return models.Teacher{}, err
	}
	_, err = db.Exec("update teachers set first_name = ? , last_name = ? , email = ? , class = ? , subject = ? where id = ? ",
		updatedTeachers.FirstName, updatedTeachers.LastName, updatedTeachers.Email, updatedTeachers.Class, updatedTeachers.Subject, id)

	if err != nil {
		fmt.Printf("Error executing update query: %v\n", err)
		return models.Teacher{}, err
	}

	updatedTeachers.ID = id // Ensure ID is in the response
	fmt.Printf("Successfully updated teacher with ID %d\n", id)
	return updatedTeachers, nil
}

func isValidSortOrder(order string) bool {
	return (order == "asc" || order == "desc")
}
func isValidFiled(filed string, params map[string]string) bool {
	_, ok := params[filed]
	return ok
}

func AddFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
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
