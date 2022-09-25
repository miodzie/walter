package log

import (
	"fmt"
	"runtime"
)

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
	fmt.Println(Gray+"[TRACE] "+Reset, args)
}

func (l ConsoleLogger) Debug(args ...interface{}) {
	fmt.Println(White+"[DEBUG] "+Reset, args)
}

func (l ConsoleLogger) Info(args ...interface{}) {
	fmt.Println(Green+"[INFO] "+Reset, args)
}

func (l ConsoleLogger) Warn(args ...interface{}) {
	fmt.Println(Yellow+"[WARN] "+Reset, args)
}

func (l ConsoleLogger) Error(args ...interface{}) {
	fmt.Println(Red+"[ERR] "+Reset, args)
}
