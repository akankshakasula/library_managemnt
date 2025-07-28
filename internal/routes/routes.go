package routes

import (
	"github.com/gofiber/fiber/v2"
	"library-management/internal/handlers"
	"library-management/internal/middleware" // Import the new middleware package
	"library-management/internal/models"     // Import models for roles
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// User Authentication Routes (No JWT required for these, they issue the JWT)
	api.Post("/signup", handlers.SignUp)
	api.Post("/signin", handlers.SignIn)

	// Protected Routes (Require Authentication)
	protected := api.Group("/").Use(middleware.Authenticate()) // Apply Auth middleware to all routes in this group

	// Book Management Routes
	protected.Get("/books", handlers.GetAllBooks) // Publicly accessible but we will put Auth for now
	protected.Post("/books", middleware.Authorize(models.RoleLibrarian), handlers.CreateBook) // Only librarians can create
	protected.Post("/books/donate", handlers.DonateBook) // Anyone logged in can donate

	// Borrowing & Returning Routes
	protected.Post("/books/borrow", handlers.BorrowBook)
	protected.Post("/books/return/:id", handlers.ReturnBook)

	// Example of a librarian-only route
	// protected.Put("/users/:id/block", middleware.Authorize(models.RoleLibrarian), handlers.BlockUser) // (Future route)
}