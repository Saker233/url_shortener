package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func InitServer() {
	r := gin.Default()
	r.GET("/health", GetHealth)


	_ = r.Run(":8080")
}

func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Healthy",
	})
}