package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	stdLogger "log"

	"github.com/oleg-polivannyi/synth-log-test/config"
	"github.com/oleg-polivannyi/synth-log-test/log"
)

func main() {
	cfg := config.LoadConfig()

	logger, err := log.NewLogger(&cfg)
	if err != nil {
		stdLogger.Fatal("Could not initialize logger:", err)
	}

	http.HandleFunc("/", handleRequest(logger, cfg))

	go sendRequestsPeriodically(logger, cfg)

	logger.Info("Starting server on port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logger.Error("Server failed:", err)
	}
}

func handleRequest(logger *log.Logger, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		logger.Info("Received request from", r.RemoteAddr, "message:", r.Form.Get("message"))
		fmt.Fprintf(w, "Hello from %s!", cfg.Tag)
	}
}

func sendRequestsPeriodically(logger *log.Logger, cfg config.Config) {
	ticker := time.NewTicker(time.Duration(60/cfg.EventFrequency) * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		message := generateGUID()
		form := fmt.Sprintf("message=%s", message)
		resp, err := http.Post(cfg.TargetURL, "application/x-www-form-urlencoded", bytes.NewBufferString(form))
		if err != nil {
			logger.Error("Failed to send request:", err)
		} else {
			logger.Info("Sent request to", cfg.TargetURL, "with response status:", resp.Status)
			resp.Body.Close()
		}
	}
}

func generateGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
