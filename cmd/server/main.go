package main

import (
	"net/http"

	"github.com/tukhvatullinsm/golang-project/internal/handlers"
	"github.com/tukhvatullinsm/golang-project/internal/storage"
)

const (
	IP   string = "localhost"
	PORT string = "8080"
)

func main() {
	// TODO: init storage object
	objStorage := storage.New()

	//TODO: init App object
	webapp := handlers.WebApp{objStorage}

	// TODO: init new handler
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", webapp.ServeHTTP)

	// TODO: Run and Check Server
	err := http.ListenAndServe(IP+":"+PORT, mux)
	if err != nil {
		panic(err)
	}
}
