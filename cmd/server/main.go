package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/tukhvatullinsm/golang-project/internal/handlers/server"
	"github.com/tukhvatullinsm/golang-project/internal/storage"
)

type Config struct {
	Endpoint string `env:"ADDRESS"`
}

func main() {

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Endpoint == "" {
		flag.StringVar(&cfg.Endpoint, "a", "localhost:8080", "Enter endpoint socket (address:port)")
		flag.Parse()
	}

	objStorage := storage.New()

	webapp := server.WebApp{}
	webapp.Init(objStorage)

	router := chi.NewRouter()
	router.Get("/value/{type}/{name}", webapp.GetValue)
	router.Get("/", webapp.GetParam)
	router.Post("/update/{type}/{name}/{value}", webapp.SetValues)

	err = http.ListenAndServe(cfg.Endpoint, router)
	if err != nil {
		log.Fatal(err)
	}
}
