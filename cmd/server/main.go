package main

import (
	"net/http"
	"strings"

	"github.com/tukhvatullinsm/golang-project/internal/handlers"
	"github.com/tukhvatullinsm/golang-project/internal/storage"
)

const (
	IP   string = "localhost"
	PORT string = "8080"
)

type AppObject struct {
	objStorage *storage.MemStorage
}

func (app *AppObject) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rawQueryParams := r.URL.EscapedPath()
		strippedQueryParams := strings.TrimPrefix(rawQueryParams, "/update/")
		sliceQueryParams := strings.Split(strippedQueryParams, "/")
		ct := r.Header.Get("Content-Type")
		switch ct {
		case "text/plain":
			res, detail := handlers.CheckUrlParams(sliceQueryParams)
			if !res {
				w.WriteHeader(detail)
				return
			}
			app.objStorage.Set(sliceQueryParams[0], sliceQueryParams[1], sliceQueryParams[2])
			w.WriteHeader(http.StatusOK)
			return
		default:
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
	}
	if r.Method == "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	// TODO: init storage object
	objStorage := storage.New()

	//TODO: init App object
	app := AppObject{objStorage}

	// TODO: init new handler
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", app.ServeHTTP)

	// TODO: Run and Check Server
	err := http.ListenAndServe(IP+":"+PORT, mux)
	if err != nil {
		panic(err)
	}
}
