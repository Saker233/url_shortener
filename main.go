package main

import (
	"url_shortener/internal/handler"
	"url_shortener/internal/util"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	godotenv.Load("app.env")

	util.ConnectToDB()
	handler.InitServer()

}
