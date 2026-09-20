package main

import (
	"url_shortener/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	godotenv.Load("app.env")
	r := gin.Default()

	r.GET("/health", handler.GetHealth)

	_ = r.Run(":8080")
}