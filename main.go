package main

import (
	stdLogger "log"

	"bytes"
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/oleg-polivannyi/synth-log-test/config"
	"github.com/oleg-polivannyi/synth-log-test/log"
)

func main() {
	config := config.LoadConfig()

	logger, err := log.NewLogger(&config)
	if err != nil {
		stdLogger.Fatal("Could not initialize logger:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the body of the request
		r.ParseForm()
		// Log the request
		logger.Info("Received request from", r.RemoteAddr, "message: ", r.Form.Get("message"))
		fmt.Fprintf(w, "Hello from %s!", config.Tag)
	})

	go func() {
		for {
			time.Sleep(time.Duration(60/config.EventFrequency) * time.Second)
			message := GUID()
			form := fmt.Sprintf("message=%s", message)
			resp, err := http.Post(config.TargetURL, "application/x-www-form-urlencoded", bytes.NewBufferString(form))
			if err != nil {
				logger.Error("Failed to send request:", err)
			} else {
				logger.Info("Sent request to", config.TargetURL, "with response status:", resp.Status)
				resp.Body.Close()
			}
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error("Server failed:", err)
	}
}

func GUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
