package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	db_host := os.Getenv("HOST")
	db_pass := os.Getenv("DB_PASSWORD")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_pass, db_host, db_port, db_name)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	
	// Set connection pool limits (optional but recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	fmt.Println("Database connection pool initialized")
	return nil
}
