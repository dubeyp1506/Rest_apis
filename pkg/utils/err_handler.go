package utils

import (
	"fmt"
	"log"
	"os"
)

var errorLogger = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)

func ErrorHandler(err error, message string) error {
	errorLogger.Output(2, fmt.Sprintf("%s: %v", message, err))
	return fmt.Errorf("%s", message)
}
