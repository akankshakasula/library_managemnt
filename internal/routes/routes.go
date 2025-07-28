package routes

import (
	"github.com/gofiber/fiber/v2"
	"library-management/internal/handlers"
	"library-management/internal/middleware" 
	"library-management/internal/models"     
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/signup", handlers.SignUp)
	api.Post("/signin", handlers.SignIn)

	protected := api.Group("/").Use(middleware.Authenticate()) 

	protected.Get("/books", handlers.GetAllBooks)
	protected.Post("/books", middleware.Authorize(models.RoleLibrarian), handlers.CreateBook)
	protected.Post("/books/donate", handlers.DonateBook) 

	protected.Post("/books/borrow", handlers.BorrowBook)
	protected.Post("/books/return/:id", handlers.ReturnBook)

}