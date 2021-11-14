package logs

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/dumacp/go-logs/internal/rotate"
)

// var PrintFuncWarning func(format string, v ...interface{})
// var PrintFuncError func(format string, v ...interface{})
// var PrintFuncInfo func(format string, v ...interface{})

//LogError log error
var LogError = New(os.Stdout, "[ ERROR ] ", log.Ldate|log.Ltime)

//LogWarn log Warning
var LogWarn = New(os.Stdout, "[ WARN ] ", log.Ldate|log.Ltime)

//LogInfo log Info
var LogInfo = New(os.Stdout, "[ INFO ] ", log.Ldate|log.Ltime)

//LogBuild log Debug
var LogBuild = New(os.Stdout, "[ BUILD ] ", log.Ldate|log.Ltime)

//Logger struct to logger
type Logger struct {
	*log.Logger

	// FuncPrintf func(format string, v ...interface{})
}

// func (l *Logger) Printf(format string, v ...interface{}) {
// 	if l.FuncPrintf != nil {
// 		l.Printf(format, v...)
// 	}
// 	l.Logger.Printf(format, v...)
// }

//New create Logger
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{log.New(out, prefix, flag)}
}

//SetLogError set logs with ERROR level
func (logg *Logger) SetLogError(logger *log.Logger) {
	logg.Logger = logger
}

//Disable set logs with ERROR level
func (logg *Logger) Disable() {
	logg.Logger.SetOutput(ioutil.Discard)
}

func NewRotate(dir, prefixname string, size int64, count, logFlat int) (*log.Logger, error) {

	return rotate.NewLogger(dir, prefixname, size, count, logFlat)

}

func NewRotateWithFuncWriter(funcWrite func([]byte) []byte, dir, prefixname string, size int64, count, logFlat int) (*log.Logger, error) {
	return rotate.NewLoggerWithFuncWriter(funcWrite, dir, prefixname, size, count, logFlat)

}
