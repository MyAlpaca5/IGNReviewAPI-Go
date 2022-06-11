package api

import (
	"fmt"
	"os"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/router"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/utils"
	"github.com/spf13/viper"
)

func Run() {
	// TODO: move this function into other file out of api folder
	// initialize viper for configuration
	utils.InitConfig()

	// initialize conn variable for database connection
	_, close, err := db.NewPGXPool()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal - %v", err)
		os.Exit(1)
	}
	defer close()

	// initialize router
	router := router.NewRouter()
	router.Run(":" + viper.GetString("port"))
}
