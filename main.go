package main

import (

	"github.com/gin-gonic/gin"
	"url_shortener/internal/handler"
)

func main() {
	r := gin.Default()

	r.GET("/health", handler.GetHealth)

	_ = r.Run(":8080")
}