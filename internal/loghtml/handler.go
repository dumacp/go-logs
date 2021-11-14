package loghtml

import (
	"html/template"
	"net/http"
	"regexp"
)

func (l *Log) GetLog(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	if l.Data == nil {
		l.Data = make(map[string]interface{})
	}

	data := l.Data["0"]
	if len(path) > 1 {
		re, err := regexp.Compile("log_([0-9]).html")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := re.FindStringSubmatch(path)
		if len(res) < 2 || len(res[1]) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		key := res[1]
		data = l.Data[key]
	}

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}
