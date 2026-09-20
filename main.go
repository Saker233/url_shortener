package main

import (
	"context"
	"fmt"
	"os"
	"url_shortener/internal/handler"

	"github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

var conn *pgx.Conn

func main() {
	godotenv.Load("app.env")

	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DB_SOURCE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unalbe to connect to database: %v\n", err)
		os.Exit(1)
	}

	r := gin.Default()

	r.GET("/health", handler.GetHealth)

	_ = r.Run(":8080")
}