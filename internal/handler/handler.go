package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func InitServer() {
	r := gin.Default()
	r.GET("/health", GetHealth)

	_ = r.Run(os.Getenv("PORT"))
}

func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Healthy",
	})
}
