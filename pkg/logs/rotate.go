package logs

import (
	"log"

	"github.com/dumacp/go-logs/internal/rotate"
)

type Rotate interface {
	Write(output []byte) (int, error)
	Close() error
	Rotate() error
	NewLogger(prefix string, logFlag int) *log.Logger
	NewLoggerWithFuncOutput(logFlag int, funcOut func(output []byte) []byte) *log.Logger
}

func NewRotate(dir, prefixname string, size int64, count int) (Rotate, error) {

	conf := &rotate.Conf{
		Directory: dir,
		Prefix:    prefixname,
		Size:      size,
		Count:     count,
	}

	return rotate.New(conf)

}

// func NewRotateWithFuncWriter(funcWrite func([]byte) []byte,
// 	dir, prefixname string, size int64, count, logFlat int) (Rotate, error) {
// 	conf := &rotate.Conf{
// 		Directory: dir,
// 		Prefix:    prefixname,
// 		Size:      size,
// 		Count:     count,
// 	}

// 	return rotate.NewWithFuncWriter(funcWrite, conf)

// }
