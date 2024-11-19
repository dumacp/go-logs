package logs

import (
	"fmt"
	"io"
	"log"
	"os"
)

// LogError log error
var LogError = New(os.Stdout, "[ ERROR ] ", log.Ldate|log.Ltime)

// LogWarn log Warning
var LogWarn = New(os.Stdout, "[ WARN ] ", log.Ldate|log.Ltime)

// LogInfo log Info
var LogInfo = New(os.Stdout, "[ INFO ] ", log.Ldate|log.Ltime)

// LogBuild log Debug
var LogBuild = New(os.Stdout, "[ BUILD ] ", log.Ldate|log.Ltime)

// Logger struct to logger
type Logger struct {
	*log.Logger
	prefix        string
	printToStdout bool
}

// New create Logger
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{log.New(out, prefix, flag), prefix, false}
}

// SetLogError set logs with ERROR level
func (logg *Logger) SetLogError(logger *log.Logger) {
	logg.Logger = logger
}

// Disable set logs with ERROR level
func (logg *Logger) Disable() {
	logg.Logger.SetOutput(io.Discard)
}

// EnableStdout enables printing to stdout
func (logg *Logger) EnableStdout() {
	logg.printToStdout = true
}

// Printf prints the log message and optionally to stdout
func (logg *Logger) Printf(format string, v ...interface{}) {
	logg.Logger.Printf(format, v...)
	if logg.printToStdout {
		if len(logg.prefix) > 0 {
			fmt.Print(logg.prefix)
		}
		fmt.Printf(format, v...)
		fmt.Println()
	}
}

// Printf prints the log message and optionally to stdout
func (logg *Logger) Print(v ...interface{}) {
	logg.Logger.Print(v...)
	if logg.printToStdout {
		if len(logg.prefix) > 0 {
			fmt.Print(logg.prefix)
		}
		fmt.Print(v...)
	}
}

// Printf prints the log message and optionally to stdout
func (logg *Logger) Println(v ...interface{}) {
	logg.Logger.Println(v...)
	if logg.printToStdout {
		if len(logg.prefix) > 0 {
			fmt.Print(logg.prefix)
		}
		fmt.Println(v...)
	}
}
