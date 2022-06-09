package controllers

import (
	"net/http"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/schemas"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// HealthcheckHandler handles "GET /api/healthcheck" endpoint. TODO: for now, just return plain text.
func HealthcheckHandler(c *gin.Context) {
	hc := schemas.Healthcheck{
		Status: "available",
		SystemInfo: schemas.SystemInfo{
			Env:     viper.GetString("env"),
			Version: viper.GetString("version"),
		},
	}
	c.JSON(http.StatusOK, hc)
}
