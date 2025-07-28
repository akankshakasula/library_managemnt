package main

import (
	"log"

	"library-management/internal/db"
	"library-management/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Middleware for Cross-Origin Resource Sharing
)

func main() {
	db.ConnectDatabase()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
