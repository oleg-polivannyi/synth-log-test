package log

import (
	"fmt"
	"io"
	"log"
	"os"
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

	return &Logger{
		config:     config,
		fileLogger: log.New(file, "", log.LstdFlags),
		stdLogger:  log.New(stdOut, "", log.LstdFlags),
	}, nil
}

func (l *Logger) Info(v ...interface{}) {
	l.fileLogger.Println(v...)
	l.stdLogger.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.fileLogger.Println(AugmentedLogger(v...))
	l.stdLogger.Println(AugmentedLogger(v...))
}

func AugmentedLogger(v ...interface{}) string {
	// time, tag, env, initial message
	time := time.Now()
	return fmt.Sprintf(
		"%s %s %s %s",
		time,
		config.LoadConfig().Tag,
		config.LoadConfig().Env,
		fmt.Sprint(v...),
	)
}
