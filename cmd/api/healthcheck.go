package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) healthcheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "status: available\nenvironment: %s\nversion: %s\nfull path:%s", app.config.env, version, c.FullPath())
}
