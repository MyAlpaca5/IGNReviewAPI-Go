package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthcheckGETHandler handles "/api/healthcheck" endpoint. TODO: for now, just return plain text.
func HealthcheckGETHandler(c *gin.Context) {
	c.String(http.StatusOK, "status: available\nfull path:%s", c.FullPath())
}
