package logger

import (
	"os"
	"sync"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger
var once sync.Once

func init() {
	once.Do(func() {
		Logger = log.New(os.Stderr)
	})
}

func Info(args ...interface{}) {
	Logger.Info(args)
}

func Warning(args ...interface{}) {
	Logger.Warn(args)
}

func Error(args ...interface{}) {
	Logger.Error(args)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args)
}
