package controllers

import (
	"net/http"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HealthcheckController struct{}

// HealthcheckHandler handles "GET /api/healthcheck" endpoint.
func (controller HealthcheckController) HealthcheckHandler(c *gin.Context) {
	hc := models.Healthcheck{
		Status: "available",
		SystemInfo: models.SystemInfo{
			Env:     viper.GetString("env"),
			Version: viper.GetString("version"),
		},
	}
	c.JSON(http.StatusOK, hc)
}
