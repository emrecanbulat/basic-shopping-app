package main

import (
	"flag"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
	"shoppingApp/internal/client"
	"shoppingApp/internal/jsonlog"
	"shoppingApp/internal/model"
	"strconv"
	"time"
)

// Define a config struct to hold all the configuration settings for our application.
type config struct {
	port int
	env  string
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {
	var cfg config
	appPort, _ := strconv.Atoi(os.Getenv("APP_PORT"))

	flag.IntVar(&cfg.port, "port", appPort, "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("APP_ENV"), "Environment (development|staging|production)")
	flag.Parse()

	client.Connections()
	model.Migrate()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})

	err := srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}
