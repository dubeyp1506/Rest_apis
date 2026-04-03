package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db_name := os.Getenv("DB_NAME")
	db, err := sqlconnect.ConnectDb(db_name)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the db", db)
	port := os.Getenv("PORT")

	rl := mw.NewRateLimiter(5, time.Minute)

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

	router := router.Router()

	SecureMx := utils.ApplyMiddleWare(router, mw.Hpp(*hppOptions), mw.Compression, mw.SecurityHandler, mw.ResponseTime, rl.RateLimiterMiddleware, mw.Cors)
	fmt.Println(SecureMx)

	server := &http.Server{
		Addr: port,
		// Handler: mw.Hpp(*hppOptions)(rl.RateLimiterMiddleware(mw.Compression(mw.ResponseTime(mw.SecurityHandler(mw.Cors(router)))))),
		Handler: mw.SecurityHandler(router),
		// Handler:   router,
		// Handler:   middlewares.Cors(router),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Listening on port", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error Starting the server :", err)
	}
}
