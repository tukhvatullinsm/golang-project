package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tukhvatullinsm/golang-project/internal/handlers/server"
	"github.com/tukhvatullinsm/golang-project/internal/storage"
)

var WebServer struct {
	Endpoint string
}

func init() {
	flag.StringVar(&WebServer.Endpoint, "a", "localhost:8080", "Enter endpoint socket (address:port)")
}

func main() {
	// TODO: Init Server (Endpoint) socket
	flag.Parse()
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
	err := http.ListenAndServe(WebServer.Endpoint, router)
	if err != nil {
		panic(err)
	}
}
