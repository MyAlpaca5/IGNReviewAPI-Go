package api

import (
	"flag"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/router"
)

const Version = "1.0.0"

// Config stores all configuration related data.
type config struct {
	port string
	env  string
}

func NewConfig() config {
	var cfg config

	// Read value from command line parameters and set them into config struct
	flag.StringVar(&cfg.port, "port", "50001", "API server port number")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|staging|prod)")
	flag.Parse()

	return cfg
}

func Run() {
	cfg := NewConfig()
	router := router.NewRouter()

	router.Run(":" + cfg.port)
}
