package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

const (
	version = "1.0.3"
)

var dirData string
var dirHtml string
var filepattern string
var showversion bool

func init() {
	flag.StringVar(&filepattern, "filepattern", "([[:alpha:]]+)_{0,1}([0-9]{1,2})\\.log$", "logs filepattern")
	flag.StringVar(&dirData, "dirData", "/SD/logs/", "logs directory")
	flag.StringVar(&dirHtml, "dirHtml", "/SD/htmllogs/", "logs directory")
	flag.BoolVar(&showversion, "version", false, "show version")
}

func main() {
	flag.Parse()

	if showversion {
		fmt.Printf("version: %s\n", version)
		os.Exit(2)
	}

	refilename, err := regexp.Compile(filepattern)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	log.Println("filedata api")
	fileserverJS := r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(fmt.Sprintf("%s/js", dirHtml)))))
	if methods, err := fileserverJS.GetMethods(); err != nil {
		for i, v := range methods {
			log.Printf("Method %d: %s", i, v)
		}
	}
	fileserverCSS := r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(fmt.Sprintf("%s/css", dirHtml)))))
	if methods, err := fileserverCSS.GetMethods(); err != nil {
		for i, v := range methods {
			log.Printf("Method %d: %s", i, v)
		}
	}

	tmpdir, err := ioutil.TempDir("", "data")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("tempdir: %s", tmpdir)
	defer os.Remove(tmpdir)
	defer func() {
		log.Printf("remove tempdir: %s", tmpdir)
		os.RemoveAll(tmpdir) // clean up
	}()

	fileserverData := r.PathPrefix("/data/").Handler(http.StripPrefix("/data/", http.FileServer(http.Dir(tmpdir))))
	if methods, err := fileserverData.GetMethods(); err != nil {
		for i, v := range methods {
			log.Printf("Method %d: %s", i, v)
		}
	}

	// directory.

	// fs := os.DirFS("/data")

	funcHandler := func(w http.ResponseWriter, r *http.Request) {
		datafiles := struct {
			Content *os.File
		}{
			Content: nil,
		}

		log.Printf("Request %+v", r)
		pathParams := mux.Vars(r)
		log.Printf("Params %+v", pathParams)
		path := pathParams["filename"]

		data := struct {
			Items []string
		}{
			Items: make([]string, 0),
		}
		key := 0
		prefix := ""
		if len(path) > 0 {
			re, err := regexp.Compile(`^([[:alpha:]]{3})[^\d]*([0-9]{1,2})\.html`)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			res := re.FindStringSubmatch(path)
			if len(res) < 3 || len(res[1]) <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			key, err = strconv.Atoi(res[2])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			prefix = res[1]
		}

		fmt.Println("regex ok")

		directory, err := os.ReadDir(dirData)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("read dir ok")
		for _, v := range directory {
			fi, err := v.Info()
			if err != nil {
				log.Println(err)
				continue
			}
			if fi.IsDir() {
				continue
			}

			fmt.Printf("parse filename: %s\n", fi.Name())
			res := refilename.FindStringSubmatch(fi.Name())
			fmt.Printf("submatch filename : %+v\n", res)
			if len(res) < 3 || len(res[1]) <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Println("regex ok")

			index, err := strconv.Atoi(res[2])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Println("03")

			fmt.Printf("prefix : %s, %s\n", prefix, fi.Name())
			if !strings.HasPrefix(fi.Name(), prefix) {
				continue
			}
			if index != key {
				continue
			}

			filename := filepath.Join(dirData, fi.Name())
			fo, err := os.Open(filename)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer fo.Close()
			datafiles.Content = fo
			break
		}

		if datafiles.Content == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// printable, err := regexp.Compile(`[\t\n\v\f\r [:graph:]]+$`)
		// printable, err := regexp.Compile(`[\t\n\v\r A-Za-z0-9!"#$%&()*+,\-./:;<=>?[\\\]^_{|}~]+$`)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		stats, _ := datafiles.Content.Stat()
		fmt.Printf("bytes in contents file: %s\n", stats)
		scan := bufio.NewScanner(datafiles.Content)

		bufferSize := 256 * 1024
		scannerBuffer := make([]byte, bufferSize)
		scan.Buffer(scannerBuffer, bufferSize)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			text := scan.Text()
			// if printable.MatchString(text) {
			itemp := make(map[string]interface{})
			if err := json.Unmarshal([]byte(text), &itemp); err != nil {
				continue
			}
			data.Items = append(data.Items, text)
			// }
		}
		// } else {
		// 	scan := bufio.NewScanner(datafiles.Contents[0])
		// 	for scan.Scan() {
		// 		data.Items = append(data.Items, scan.Text())
		// 	}
		// }
		fmt.Printf("Items count in log: %d\n", len(data.Items))

		tmpfiles, err := os.ReadDir(tmpdir)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// log.Println("0")
		for _, v := range tmpfiles {
			if v.IsDir() {
				continue
			}
			filename := filepath.Join(tmpdir, v.Name())
			os.Remove(filename)
		}

		tjs := template.New("data")

		// log.Println("1")
		tjs, err = tjs.Funcs(template.FuncMap{"StringsJoin": strings.Join}).Parse(tpjs)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tfdata, err := ioutil.TempFile(tmpdir, "data")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer tfdata.Close()
		// log.Println("2")

		err = tjs.Execute(tfdata, data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// log.Println("3")

		thtml, err := template.New("webpage").Parse(tpl)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filename := filepath.Base(tfdata.Name())
		err = thtml.Execute(w, map[string]interface{}{
			"Files": []string{filename},
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// log.Println("5")

		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}

	apiv1 := r.PathPrefix("/logs").Subrouter()
	apiv1.HandleFunc("/{filename}", funcHandler).Methods(http.MethodGet)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)

	for range finish {
		log.Print("Finish")
		break
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
