package internal

import (
	"fmt"
	"os"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/router"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/utils"
	"github.com/spf13/viper"
)

func Run() {
	// initialize viper for configuration
	utils.InitConfig()

	// create repo variable for database interaction
	pool, close, err := db.InitDBConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB ERROR - %v", err)
		os.Exit(1)
	}
	defer close()

	// create new router
	router := router.NewRouter(pool)
	router.Run(":" + viper.GetString("port"))
}
