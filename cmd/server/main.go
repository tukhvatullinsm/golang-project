package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/tukhvatullinsm/golang-project/internal/handlers/server"
	"github.com/tukhvatullinsm/golang-project/internal/storage"
	"go.uber.org/zap"
)

type Config struct {
	Endpoint string `env:"ADDRESS"`
}

func main() {
	logObj, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger := *logObj.Sugar()

	defer logger.Sync()

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Endpoint == "" {
		logger.Info(
			"Failed to get configuration from environment, continue to get configuration from CLI or default value")
		flag.StringVar(&cfg.Endpoint, "a", "localhost:8080", "Enter endpoint socket (address:port)")
		flag.Parse()
	} else {
		logger.Info(
			"Success to get configuration from environment variable")
	}

	objStorage := storage.New()

	webapp := server.WebApp{}
	webapp.Init(objStorage, &logger)

	router := chi.NewRouter()
	router.Use(webapp.LoggingMiddleware)
	router.Post("/value/", webapp.GetValueJSON)
	router.Get("/value/{type}/{name}", webapp.GetValue)
	router.Get("/", webapp.GetParam)
	router.Post("/update/", webapp.SetValuesJSON)
	router.Post("/update/{type}/{name}/{value}", webapp.SetValues)

	logger.Infow("Starting server", "addr", cfg.Endpoint)
	err = http.ListenAndServe(cfg.Endpoint, router)
	if err != nil {
		logger.Fatal("Web Server cannot run, so ", err)
	}
}
