package main

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	database.Connect()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	app := fiber.New()
	routes.Setup(app)

	if err = app.Listen(":" + port); err != nil {
		panic(err)
	}
}
