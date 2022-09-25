package log

import (
	"fmt"
	"runtime"
)

// TODO: Extract to color package.

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

// ConsoleLogger is a basic logger that prints to the console.
type ConsoleLogger struct {
}

func (l ConsoleLogger) Trace(args ...interface{}) {
	fmt.Printf(gray())
	fmt.Println(args...)
}

func (l ConsoleLogger) Debug(args ...interface{}) {
	fmt.Printf(white())
	fmt.Println(args...)
}

func (l ConsoleLogger) Info(args ...interface{}) {
	fmt.Printf(green())
	fmt.Println(args...)
}

func (l ConsoleLogger) Warn(args ...interface{}) {
	fmt.Printf(yellow())
	fmt.Println(args...)
}

func (l ConsoleLogger) Error(args ...interface{}) {
	fmt.Printf(red())
	fmt.Println(args...)
}

func (l ConsoleLogger) Tracef(format string, args ...interface{}) {
	fmt.Printf(gray()+format, args...)
}

func (l ConsoleLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf(white()+format, args...)
}

func (l ConsoleLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(green()+format, args...)
}

func (l ConsoleLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf(yellow()+format, args...)
}

func (l ConsoleLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(red()+format, args...)
}

func gray() string {
	return Gray + "[TRACE] " + Reset
}
func white() string {
	return White + "[DEBUG] " + Reset
}

func green() string {
	return Green + "[INFO] " + Reset
}

func yellow() string {
	return Yellow + "[WARN] " + Reset
}

func red() string {
	return Red + "[ERR] " + Reset
}
