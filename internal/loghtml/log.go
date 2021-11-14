package loghtml

import "time"

type Log struct {
	DirFiles string
	Data     map[string]interface{}
}

type Data struct {
	Title   string
	Type    string
	Message []byte
	Date    time.Time
}
