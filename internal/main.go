package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/router"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/utils"
	"github.com/spf13/viper"
)

func Run() {
	// --- initialize viper for configuration ---
	utils.InitConfig()

	// --- create database pool connection ---
	pool, close, err := db.NewPool()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB ERROR - %v", err)
		os.Exit(1)
	}
	defer close()

	// --- create new router ---
	router := router.New(pool)

	// --- gracefully shutdown ---
	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	srv := &http.Server{
		Addr:    ":" + viper.GetString("port"),
		Handler: router,
	}

	// create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server Error: %s\n", err.Error())
		}
	}()

	// listen on catched signals
	<-ctx.Done()

	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeout); err != nil {
		fmt.Printf("Server forced to shutdown: %s\n", err.Error())
	}

	fmt.Println("Server exiting")
}
