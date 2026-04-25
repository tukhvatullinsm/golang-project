package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/tukhvatullinsm/golang-project/internal/handlers/server"
	"github.com/tukhvatullinsm/golang-project/internal/storage"
	"go.uber.org/zap"
)

type Config struct {
	Endpoint      string `env:"ADDRESS"`
	StoreInterval *int64 `env:"STORE_INTERVAL"`
	StorePath     string `env:"FILE_STORE_PATH"`
	Restore       *bool  `env:"RESTORE"`
}

func main() {
	logObj, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger := *logObj.Sugar()
	defer logger.Sync()

	var cfg Config
	dateString := time.Now().Format("2006-01-02")
	defaultFilename := "GOProject" + dateString + ".json"
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Endpoint == "" {
		flag.StringVar(&cfg.Endpoint, "a", "localhost:8080", "Enter endpoint socket (address:port)")
	}
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = new(int64)
		flag.Int64Var(cfg.StoreInterval, "i", 300, "Enter store interval (seconds)")
	}
	if cfg.StorePath == "" {
		flag.StringVar(&cfg.StorePath, "f", defaultFilename, "Enter path to store files")
	}
	if cfg.Restore == nil {
		cfg.Restore = new(bool)
		flag.BoolVar(cfg.Restore, "r", true, "Enter state of restore data from local file or not")
	}
	flag.Parse()
	objStorage := storage.New(&logger)

	webapp := server.WebApp{}
	webapp.Init(objStorage, &logger, cfg.StorePath, *cfg.Restore, *cfg.StoreInterval)
	defer webapp.ObjStorage.Close()
	router := chi.NewRouter()
	router.Use(webapp.LoggingMiddleware)
	router.Use(webapp.GzipMiddleware)
	router.Use(webapp.PanicMiddleware)
	router.Post("/value/", webapp.GetValueJSON)
	router.Get("/value/{type}/{name}", webapp.GetValue)
	router.Get("/", webapp.GetParam)
	router.Post("/update/", webapp.SetValuesJSON)
	router.Post("/update/{type}/{name}/{value}", webapp.SetValues)

	logger.Infow("Starting server", "addr", cfg.Endpoint)
	timer := time.AfterFunc(time.Duration(0)*time.Second, webapp.SaveValue)
	defer timer.Stop()
	err = http.ListenAndServe(cfg.Endpoint, router)
	if err != nil {
		logger.Fatal("Web Server cannot run, so ", err)
	}

}
