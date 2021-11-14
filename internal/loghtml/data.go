package loghtml

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ServeFilesData(r *mux.Route, dir string) {
	fileserver := r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(dir))))
	if methods, err := fileserver.GetMethods(); err != nil {
		for i, v := range methods {
			log.Printf("Method %d: %s", i, v)
		}
	}
}
