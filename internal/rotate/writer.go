package rotate

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

type Conf struct {
	Directory string
	Prefix    string
	Size      int64
	Count     int
}

// type funcWriter func(output []byte) []byte

type RotateWriter struct {
	Conf *Conf
	lock sync.Mutex
	fp   *os.File
	// funcWrite funcWriter
}

// Make a new RotateWriter. Return nil if error occurs during setup.
func New(conf *Conf) (*RotateWriter, error) {
	if fd, err := os.Open(conf.Directory); err != nil {
		if err := os.Mkdir(conf.Directory, 0755); err != nil {
			return nil, err
		}
	} else {
		if stats, err := fd.Stat(); err != nil {
			if err := os.Mkdir(conf.Directory, 0755); err != nil {
				return nil, err
			}
		} else if !stats.IsDir() {
			return nil, fmt.Errorf("%s is not directory", conf.Directory)
		} else {
			fd.Close()
		}
	}

	filename := filepath.Join(conf.Directory, fmt.Sprintf("%s_0.log", conf.Prefix))
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	w := &RotateWriter{
		Conf: conf,
		lock: sync.Mutex{},
		fp:   fp,
	}
	return w, nil
}

// Make a new RotateWriter. Return nil if error occurs during setup.
func NewWithFuncWriter(funcWriter func(output []byte) []byte, conf *Conf) (*RotateWriter, error) {

	if fd, err := os.Open(conf.Directory); err != nil {
		if err := os.Mkdir(conf.Directory, 0755); err != nil {
			return nil, err
		}
	} else {
		if stats, err := fd.Stat(); err != nil {
			if err := os.Mkdir(conf.Directory, 0755); err != nil {
				return nil, err
			}
		} else if !stats.IsDir() {
			return nil, fmt.Errorf("%s is not directory", conf.Directory)
		} else {
			fd.Close()
		}
	}

	filename := filepath.Join(conf.Directory, fmt.Sprintf("%s_0.log", conf.Prefix))
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	w := &RotateWriter{
		Conf: conf,
		lock: sync.Mutex{},
		fp:   fp,
		// funcWrite: funcWriter,
	}
	return w, nil
}

func (w *RotateWriter) NewLogger(prefix string, logFlag int) *log.Logger {
	return log.New(w, prefix, logFlag)
}

type bypass struct {
	out        io.Writer
	funcOutput func(output []byte) []byte
}

func (b *bypass) Write(data []byte) (int, error) {
	result := b.funcOutput(data)
	return b.out.Write(result)
}

func (w *RotateWriter) NewLoggerWithFuncOutput(logFlag int,
	funcOut func(output []byte) []byte) *log.Logger {

	byp := &bypass{
		out:        w,
		funcOutput: funcOut,
	}
	return log.New(byp, "", logFlag)
}

// Close close file.
func (w *RotateWriter) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.Close()
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	stats, err := w.fp.Stat()
	if err != nil {
		return 0, err
	}
	if stats.Size()+int64(len(output)) > w.Conf.Size {
		if err := w.Rotate(); err != nil {
			return 0, err
		}
	}
	// // log.Printf("escribiendo: %s", output)
	// if w.funcWrite != nil {
	// 	output = w.funcWrite(output)
	// }
	return w.fp.Write(output)
}

// Rotate files in directory.
func (w *RotateWriter) Rotate() error {
	// w.lock.Lock()
	// defer w.lock.Unlock()
	// Close existing file if open
	if w.fp != nil {
		if err := w.fp.Close(); err != nil {
			w.fp = nil
			return err
		}
	}
	// Rename dest file if it already exists
	dirs, err := os.ReadDir(w.Conf.Directory)
	if err != nil {
		return err
	}
	filenames := make(map[int]string)
	for _, f := range dirs {
		if f.IsDir() {
			continue
		}
		re, err := regexp.Compile(fmt.Sprintf("%s_([0-%d]).log", w.Conf.Prefix, w.Conf.Count))
		if err != nil {
			return err
		}
		res := re.FindStringSubmatch(f.Name())
		if len(res) < 2 || len(res[1]) <= 0 {
			continue
		}
		key, _ := strconv.Atoi(res[1])
		filenames[key] = f.Name()
	}
	keys := make([]int, 0)
	for k, _ := range filenames {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	log.Printf("keys: %v", keys)
	for len(keys) > 0 && keys[0] >= w.Conf.Count {
		keys = keys[1:]
	}
	// Move files
	for _, k := range keys {
		newname := filepath.Join(w.Conf.Directory, fmt.Sprintf("%s_%d.log", w.Conf.Prefix, k+1))
		oldname := filepath.Join(w.Conf.Directory, filenames[k])
		if err := os.Rename(oldname, newname); err != nil {
			return err
		}
	}
	// Create a new first file.
	newfilename := filepath.Join(w.Conf.Directory, fmt.Sprintf("%s_%d.log", w.Conf.Prefix, 0))
	w.fp, err = os.Create(newfilename)
	if err != nil {
		return err
	}
	return nil
}
