package seras

import (
	"io"
	"log"
	"os"
	"time"
)

var logger *log.Logger

func InitLogger() {
	logger = NewLogger(GetDefaultOptions())
}

func Log(msg string) {
	logger.Println(msg)
}

type Options struct {
	TimeFormat   string
	LogDirectory string
}

func GetDefaultOptions() Options {
	options := Options{}
	options.TimeFormat = "2006-Jan-2-Mon-15_04_05"
	options.LogDirectory = "~/.tsugumi/logs"
	return options
}

func NewLogger(options Options) *log.Logger {
	if _, err := os.Stat(options.LogDirectory); os.IsNotExist(err) {
		os.MkdirAll(options.LogDirectory, os.ModePerm)
	}
	logFile, err := os.OpenFile(options.LogDirectory+"/"+time.Now().Format(options.TimeFormat), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	return log.New(multiWriter, "", log.Ldate|log.Ltime)
}
