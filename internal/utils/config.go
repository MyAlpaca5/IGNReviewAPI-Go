package utils

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitConfig() {
	// Read value from command line parameters and set them into config struct
	flag.String("port", "50001", "API server port number")
	flag.String("env", "dev", "Environment (dev|staging|prod)")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	viper.SetDefault("version", "1.0.0")
	viper.BindEnv("db_url", "IGN_DB_URL")
}
