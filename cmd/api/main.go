package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// Config stores all configuration related data.
type config struct {
	port int
	env  string
}

// Application stores the dependencies for HTTP handlers, helpers, and middlewares.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	// Read value from command line parameters and set them into config struct
	flag.IntVar(&cfg.port, "port", 4000, "API server port number")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|staging|prod)")
	flag.Parse()

	// Write message to stdout prefixing with date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// app is a application specific varialbe to store related data
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Initialize servemux
	mux := http.NewServeMux()

	// Bind handlers
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// Initialize HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  2 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start HTTP server
	logger.Printf("API %s server is listening on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
