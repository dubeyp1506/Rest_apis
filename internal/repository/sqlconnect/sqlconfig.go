package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb(dbname string) (*sql.DB, error) {
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	db_host := os.Getenv("HOST")
	db_pass := os.Getenv("DB_PASSWORD")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_pass, db_host, db_port, db_name)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to db")
	return db, nil
}
