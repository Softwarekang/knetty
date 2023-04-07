/*
	Copyright 2022 Phoenix

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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
	log.Printf(format, args...)
}

func (l defaultLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (l defaultLogger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func (l defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l defaultLogger) Info(args ...interface{}) {
	log.Print(args...)
}

func (l defaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l defaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l defaultLogger) Debug(args ...interface{}) {
	log.Print(args...)
}
