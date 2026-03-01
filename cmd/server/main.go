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

func init() {

}

func main() {
	// TODO: Init Server Configuration (Endpoint)
	//os.Setenv("ADDRESS", "localhost:100")

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Endpoint == "" {
		flag.StringVar(&cfg.Endpoint, "a", "localhost:8080", "Enter endpoint socket (address:port)")
		flag.Parse()
	}

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
	err = http.ListenAndServe(cfg.Endpoint, router)
	if err != nil {
		panic(err)
	}
}
