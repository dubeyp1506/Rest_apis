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

	if err := sqlconnect.InitDB(); err != nil {
		panic(err)
	}
	fmt.Println("Connected to the db", sqlconnect.DB)
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

	router := router.MainRouter()

	// Wrap the JWTMiddleware to exclude specific routes like login
	jwtAuth := mw.ExcludeMiddleware(
		mw.JWTMiddleware,
		"/login",           // Adjust these paths to precisely match your actual login route
		"/forgot-password", // e.g. "/execs/login" or similar
	)

	// Apply all middlewares, including our new jwtAuth
	SecureMx := utils.ApplyMiddleWare(
		router,
		mw.Hpp(*hppOptions),
		mw.Compression,
		mw.SecurityHandler,
		mw.ResponseTime,
		rl.RateLimiterMiddleware,
		jwtAuth,
		mw.Cors,
	)
	fmt.Println(SecureMx)

	server := &http.Server{
		Addr: port,
		// Set the server to use the fully secured multiplexer
		Handler:   SecureMx,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Listening on port", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error Starting the server :", err)
	}
}
