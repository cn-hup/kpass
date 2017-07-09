package logger

import (
	"os"

	"github.com/teambition/gear"
	"github.com/teambition/gear/logging"
)

// std is global logging instance ...
var std = logging.New(os.Stdout)

// Default returns the global logging instance
func Default() *logging.Logger {
	return std
}

// Init ...
func Init() {
	std.SetLevel(logging.InfoLevel)
}

// FromCtx ...
func FromCtx(ctx *gear.Context) logging.Log {
	return std.FromCtx(ctx)
}

// Info ...
func Info(v interface{}) {
	std.Info(v)
}

// Err ...
func Err(v interface{}) {
	std.Err(v)
}

// Fatal ...
func Fatal(v interface{}) {
	std.Fatal(v)
}

// Println ...
func Println(args ...interface{}) {
	std.Println(args...)
}

// Printf ...
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}
