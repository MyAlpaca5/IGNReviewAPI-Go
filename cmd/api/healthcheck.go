package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("status: available")
	app.logger.Printf("environment: %s\n", app.config.env)
	app.logger.Printf("version: %s\n", version)
}
