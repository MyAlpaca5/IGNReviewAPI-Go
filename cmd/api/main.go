package main

import (
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

const version = "1.0.0"

// Config stores all configuration related data.
type config struct {
	port string
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
	flag.StringVar(&cfg.port, "port", "50001", "API server port number")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|staging|prod)")
	flag.Parse()

	// Write message to stdout prefixing with date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// app is a application specific varialbe to store related data
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// define routes and handlers
	v1 := r.Group("/v1")
	{
		v1.GET("/healthcheck", app.healthcheckHandler)
	}

	// run server
	r.Run(":" + app.config.port)
}
