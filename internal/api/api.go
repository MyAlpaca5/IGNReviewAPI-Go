package api

import (
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/router"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/utils"
	"github.com/spf13/viper"
)

func Run() {
	utils.InitConfig()
	router := router.NewRouter()

	router.Run(":" + viper.GetString("port"))
}
