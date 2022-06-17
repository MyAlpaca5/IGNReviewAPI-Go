package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HealthcheckController struct{}

type SystemInfo struct {
	Env     string `json:"env"`
	Version string `json:"version"`
}

type Healthcheck struct {
	Status     string     `json:"status"`
	SystemInfo SystemInfo `json:"system_info"`
}

// HealthcheckHandler handles "GET /api/healthcheck" endpoint.
func (ctrl HealthcheckController) HealthcheckHandler(c *gin.Context) {
	hc := Healthcheck{
		Status: "available",
		SystemInfo: SystemInfo{
			Env:     viper.GetString("env"),
			Version: viper.GetString("general.version"),
		},
	}
	c.JSON(http.StatusOK, hc)
}
