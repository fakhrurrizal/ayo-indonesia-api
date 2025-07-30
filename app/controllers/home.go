package controllers

import (
	"ayo-indonesia-api/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome " + config.LoadConfig().AppName,
	})
}
