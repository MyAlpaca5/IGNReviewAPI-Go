package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HealthcheckController struct{}

type systemInfo struct {
	Env     string `json:"env"`
	Version string `json:"version"`
}

type healthcheck struct {
	Status     string     `json:"status"`
	SystemInfo systemInfo `json:"system_info"`
}

// HealthcheckHandler handles "GET /healthcheck" endpoint.
func (ctrl HealthcheckController) HealthcheckHandler(c *gin.Context) {
	hc := healthcheck{
		Status: "available",
		SystemInfo: systemInfo{
			Env:     viper.GetString("env"),
			Version: viper.GetString("general.version"),
		},
	}
	c.JSON(http.StatusOK, hc)
}
