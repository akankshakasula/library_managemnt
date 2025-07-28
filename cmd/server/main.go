package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Middleware for Cross-Origin Resource Sharing
	"library-management/internal/db"
	"library-management/internal/routes"
)

func main() {
	db.ConnectDatabase()

	app := fiber.New()

	// 3. Enable CORS middleware. This is important if you'll have a frontend
	//    running on a different port or domain. For development, allowing all origins is common.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allows all origins (e.g., your frontend on localhost:3001 can talk to backend on :3000)
		AllowHeaders: "Origin, Content-Type, Accept", // Specify allowed headers
	}))

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}