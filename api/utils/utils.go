// utils/utils.go
package utils

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

// LogResponse logs the response to a file.
func LogResponse(responseBody string, endpoint string) {
	logFile, err := os.OpenFile("test/screenshots/api_responses.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Printf("[%s] Response from %s: %s\n", time.Now().Format(time.RFC3339), endpoint, responseBody)
}

// CaptureResponse saves the response body for debugging.
func CaptureResponse(responseBody string, statusCode int) {
	fileName := "test/screenshots/response_" + time.Now().Format("20060102_150405") + ".txt"
	err := ioutil.WriteFile(fileName, []byte(responseBody), 0644)
	if err != nil {
		log.Printf("Failed to capture response: %v", err)
	} else {
		log.Printf("Captured response at %s with status code %d", fileName, statusCode)
	}
}
