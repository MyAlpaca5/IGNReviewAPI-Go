package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "status: available\nfull path:%s", c.FullPath())
}
