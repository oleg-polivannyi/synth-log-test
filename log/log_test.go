package log

import (
	"os"
	"testing"

	"github.com/oleg-polivannyi/synth-log-test/config"
)

func TestNewLogger(t *testing.T) {
	cfg := config.Config{
		FileName: "test.log",
		StdOut:   true,
	}

	logger, err := NewLogger(&cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	if logger.fileLogger == nil {
		t.Fatal("Expected fileLogger to be non-nil")
	}

	if logger.stdLogger == nil {
		t.Fatal("Expected stdLogger to be non-nil")
	}

	// Clean up
	os.Remove("test.log")
}

func TestLogger_Info(t *testing.T) {
	cfg := config.Config{
		FileName: "test.log",
		StdOut:   true,
	}

	logger, err := NewLogger(&cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	logger.Info("This is an info message")

	// Check if the log file contains the message
	file, err := os.Open("test.log")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if stat.Size() == 0 {
		t.Fatal("Expected log file to contain the message")
	}

	// Clean up
	os.Remove("test.log")
}

func TestLogger_Error(t *testing.T) {
	cfg := config.Config{
		FileName: "test.log",
		StdOut:   true,
	}

	logger, err := NewLogger(&cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	logger.Error("This is an error message")

	// Check if the log file contains the message
	file, err := os.Open("test.log")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if stat.Size() == 0 {
		t.Fatal("Expected log file to contain the message")
	}

	// Clean up
	os.Remove("test.log")
}
