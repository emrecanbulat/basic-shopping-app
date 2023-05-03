package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
	"shoppingApp/internal/client"
	"shoppingApp/internal/jsonlog"
	"shoppingApp/internal/model"
	"shoppingApp/internal/seed"
	"strconv"
	"time"
)

var EnvFileError = errors.New("error loading .env file. Please make sure you are working in the correct directory")

// Define a config struct to hold all the configuration settings for our application.
type config struct {
	port int
	env  string
	jwt  struct {
		secret string
	}
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config config
	logger *jsonlog.Logger
}

func Hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Hello From Shopping App")
}

func main() {
	var cfg config
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Load environment file
	err := godotenv.Load(".env")
	if err != nil {
		logger.PrintFatal(EnvFileError, nil)
	}

	appPort, _ := strconv.Atoi(os.Getenv("APP_PORT"))

	flag.IntVar(&cfg.port, "port", appPort, "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("APP_ENV"), "Environment (development|staging|production)")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT secret")
	flag.Parse()

	client.Connections() // database connection
	model.Migrate()      // database migration
	seed.Seed()          // seed dummy data

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

	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}
