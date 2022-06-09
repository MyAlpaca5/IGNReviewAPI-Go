package api

import (
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/config"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/router"
	"github.com/spf13/viper"
)

func Run() {
	config.Setup()
	router := router.NewRouter()

	router.Run(":" + viper.GetString("port"))
}
