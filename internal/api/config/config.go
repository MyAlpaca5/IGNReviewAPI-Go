package config

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const Version = "1.0.0"

func Setup() {
	// Read value from command line parameters and set them into config struct
	flag.String("port", "50001", "API server port number")
	flag.String("env", "dev", "Environment (dev|staging|prod)")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	// hard-code config
	viper.SetDefault("version", "1.0.0")
}
