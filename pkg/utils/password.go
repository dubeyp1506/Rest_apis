package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"restapi/internal/models"
	"strings"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(user models.Exec, req models.Exec) bool {
	password := strings.Split(user.Password, ".")
	if len(password) != 2 {
		ErrorHandler(errors.New("invalid envcoded hash format"), "Internal Server Error")
		return false
	}
	saltBase64 := password[0]
	hashBase64 := password[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		ErrorHandler(errors.New("Failed to decode salt"), "Internal Server Error")
		return false
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		ErrorHandler(errors.New("Failed to encode hash password"), "Internal Server Error")
		return false
	}

	hash := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 2, 32)

	if len(hash) != len(hashedPassword) {
		ErrorHandler(errors.New("invalid Password"), "Invalid Password")
		return false
	}

	if subtle.ConstantTimeCompare(hash, hashedPassword) != 1 {
		ErrorHandler(errors.New("invalid Password"), "Invalid Password")
		return false
	}
	return true
}
