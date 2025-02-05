package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/oleg-polivannyi/synth-log-test/config"
)

type Logger struct {
	config     *config.Config
	fileLogger *log.Logger
	stdLogger  *log.Logger
}

func NewLogger(config *config.Config) (*Logger, error) {
	var file io.Writer
	var err error

	if config.FileName != "" {
		file, err = os.OpenFile(config.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	} else {
		file = io.Discard
	}

	var stdOut io.Writer
	if config.StdOut {
		stdOut = os.Stdout
	} else {
		stdOut = io.Discard
	}

	logger := &Logger{
		config:     config,
		fileLogger: log.New(file, "", 0),
		stdLogger:  log.New(stdOut, "", 0),
	}

	return logger, nil
}

func (l *Logger) Info(v ...interface{}) {
	l.fileLogger.Println(AugmentedLogger(v...))
	l.stdLogger.Println(AugmentedLogger(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.fileLogger.Println(AugmentedLogger(v...))
	l.stdLogger.Println(AugmentedLogger(v...))
}

func AugmentedLogger(v ...interface{}) string {
	// time, tag, env, initial messagetime := time.Now()
	templateStr := `level===INFO	time==={{.Time}}	path==={{.Path}}	service_name==={{.ServiceName}}	trace_id==={{.TraceID}}	request_id==={{.RequestID}}	lockbox_id==={{.LockboxID}}	request_url==={{.RequestURL}}	iam_token==={{.IAMToken}}`
	tmpl, err := template.New("log").Parse(templateStr)

	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Time        string
		Path        string
		ServiceName string
		TraceID     string
		RequestID   string
		LockboxID   string
		RequestURL  string
		IAMToken    string
	}{
		Time:        time.Now().Format("2006-01-02 15:04:05,000"),
		Path:        "path",
		ServiceName: "log-synth-test",
		TraceID:     "trace_id",
		RequestID:   "request_id",
		LockboxID:   "lockbox_id",
		RequestURL:  "request_url",
		IAMToken:    "iam_token",
	}
	var augmentedLog bytes.Buffer
	err = tmpl.Execute(&augmentedLog, data)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s %v", augmentedLog.String(), v)
}
