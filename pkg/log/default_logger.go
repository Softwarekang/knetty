package log

import (
	"log"
	"os"
)

// newDefaultLogger returns a Logger which will write log messages to stdout, and
// use same formatting runes as the stdlib log.Logger
func newDefaultLogger() Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// A defaultLogger provides a minimalistic logger satisfying the Logger interface.
type defaultLogger struct {
	logger *log.Logger
}

func (l defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (l defaultLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

func (l defaultLogger) Fatal(args ...interface{}) {
	log.Fatal(args)
}

func (l defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (l defaultLogger) Info(args ...interface{}) {
	log.Print(args)
}

func (l defaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (l defaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func (l defaultLogger) Debug(args ...interface{}) {
	log.Print(args)
}
