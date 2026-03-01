package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tukhvatullinsm/golang-project/internal/handlers/server"
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
	webapp := server.WebApp{}
	webapp.Init(objStorage)

	// TODO: init new handler
	router := chi.NewRouter()
	router.Get("/value/{type}/{name}", webapp.GetValue)
	router.Get("/", webapp.GetParam)
	router.Post("/update/{type}/{name}/{value}", webapp.SetValues)

	// TODO: Run and Check Server
	err := http.ListenAndServe(IP+":"+PORT, router)
	if err != nil {
		panic(err)
	}
}
