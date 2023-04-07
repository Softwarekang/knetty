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

var (
	DefaultLogger Logger
)

func init() {
	DefaultLogger = newDefaultLogger()
}

// Logger  A Logger is a minimalistic interface for the knetty to log messages to. Should
// be used to provide custom logging writers for the knetty to use.
type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args)
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args)
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args)
}
