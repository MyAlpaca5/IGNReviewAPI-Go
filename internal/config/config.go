package config

import (
	"fmt"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitConfig() {
	// read value from command line parameter
	env := flag.String("env", "dev", "Environment (dev|prod)")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	// load respective configurations
	if *env == "dev" {
		loadDevConfig()
	} else {
		loadProdConfig()
	}
}

func loadDevConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("couldn't load config: %s", err))
	}
}

func loadProdConfig() {
	// not implemented yet
	panic("no handler for loading production environment configuration")
}
