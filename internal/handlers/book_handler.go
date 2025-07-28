package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library-management/internal/db"
	"library-management/internal/models"
)


type DonateBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Number      string `json:"number"` 
	Genre       string `json:"genre"`
	DonatedByID uint   `json:"donated_by_id"` 
}

type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Number string `json:"number"` 
	Genre  string `json:"genre"`
}

type BorrowBookRequest struct {
	BookID uint `json:"book_id"`
	UserID uint `json:"user_id"`
}



func DonateBook(c *fiber.Ctx) error {
	req := new(DonateBookRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON body"})
	}

	if req.Title == "" || req.Author == "" || req.Number == "" || req.Genre == "" || req.DonatedByID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title, Author, Number, Genre, and DonatedByID are required for donation"})
	}

	var donorUser models.User
	if err := db.DB.First(&donorUser, req.DonatedByID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "DonatedByID does not correspond to an existing user"})
		}
		log.Printf("Database error checking donor user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	book := models.Book{
		Title:       req.Title,
		Author:      req.Author,
		Number:      req.Number,
		Genre:       req.Genre,
		DonatedByID: req.DonatedByID,
		Available:   true,
	}

	var existingBook models.Book
	if err := db.DB.Where("number = ?", book.Number).First(&existingBook).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Book with this unique number already exists"})
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Database error checking for existing book: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if err := db.DB.Create(&book).Error; err != nil {
		log.Printf("Error creating book: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not donate book"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Book donated successfully",
		"book": fiber.Map{
			"id":          book.ID,
			"title":       book.Title,
			"author":      book.Author,
			"number":      book.Number,
			"genre":       book.Genre,
			"donated_by_id": book.DonatedByID,
			"available":   book.Available,
		},
	})
}


func CreateBook(c *fiber.Ctx) error {
	req := new(CreateBookRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON body"})
	}

	if req.Title == "" || req.Author == "" || req.Number == "" || req.Genre == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title, Author, Number, and Genre are required"})
	}

	book := models.Book{
		Title:  req.Title,
		Author: req.Author,
		Number: req.Number,
		Genre:  req.Genre,
		Available: true, 
	}

	var existingBook models.Book
	if err := db.DB.Where("number = ?", book.Number).First(&existingBook).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Book with this unique number already exists"})
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Database error checking for existing book: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if err := db.DB.Create(&book).Error; err != nil {
		log.Printf("Error creating book: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create book"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Book created successfully",
		"book": fiber.Map{
			"id":        book.ID,
			"title":     book.Title,
			"author":    book.Author,
			"number":    book.Number,
			"genre":     book.Genre,
			"available": book.Available,
		},
	})
}


func GetAllBooks(c *fiber.Ctx) error {
	var books []models.Book
	if err := db.DB.Order("title ASC").Find(&books).Error; err != nil {
		log.Printf("Database error getting all books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve books"})
	}

	if len(books) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "No books found", "books": []models.Book{}})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Books retrieved successfully",
		"books":   books,
	})
}


func BorrowBook(c *fiber.Ctx) error {
	req := new(BorrowBookRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON body"})
	}

	if req.BookID == 0 || req.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BookID and UserID are required"})
	}

	var book models.Book
	if err := db.DB.Where("id = ? AND available = ?", req.BookID, true).First(&book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found or not available"})
		}
		log.Printf("Database error finding book for borrowing: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var user models.User
	if err := db.DB.Where("id = ? AND blocked = ?", req.UserID, false).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found or is blocked"})
		}
		log.Printf("Database error finding user for borrowing: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if user.Role == models.RoleStudent {
		var borrowedBooksCount int64
		db.DB.Model(&models.Borrow{}).Where("user_id = ? AND returned = ?", req.UserID, false).Count(&borrowedBooksCount)
		if borrowedBooksCount >= 3 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Student has reached the maximum borrowing limit of 3 books"})
		}
	}

	if err := db.DB.Model(&book).Update("available", false).Error; err != nil {
		log.Printf("Error updating book availability: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update book availability"})
	}

	borrowDate := time.Now()
	dueDate := borrowDate.AddDate(0, 0, 7)

	borrow := models.Borrow{
		BookID:     req.BookID,
		UserID:     req.UserID,
		BorrowDate: borrowDate,
		DueDate:    dueDate,
		Returned:   false,
	}

	if err := db.DB.Create(&borrow).Error; err != nil {
		log.Printf("Error creating borrow record: %v", err)
		db.DB.Model(&book).Update("available", true) 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to record borrow transaction"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Book borrowed successfully",
		"borrow_id": borrow.ID,
		"book_id":   borrow.BookID,
		"user_id":   borrow.UserID,
		"due_date":  borrow.DueDate,
	})
}

func ReturnBook(c *fiber.Ctx) error {
	borrowID, err := strconv.ParseUint(c.Params("id"), 10, 32) 
	if err != nil || borrowID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid borrow ID"})
	}

	var borrow models.Borrow
	if err := db.DB.Where("id = ? AND returned = ?", borrowID, false).First(&borrow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Active borrow record not found for this ID"})
		}
		log.Printf("Database error finding borrow record for return: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	returnDate := time.Now()
	fineAmount := 0.0
	if returnDate.After(borrow.DueDate) {
		overdueDuration := returnDate.Sub(borrow.DueDate)
		overdueDays := int(overdueDuration.Hours() / 24)
		if overdueDuration.Hours()/24 > float64(overdueDays) { 
			overdueDays++
		}
		fineAmount = float64(overdueDays) * 1.0
	}

	// Update borrow record
	borrow.ReturnDate = &returnDate 
	borrow.Returned = true           
	borrow.FineAmount = fineAmount   

	if err := db.DB.Save(&borrow).Error; err != nil {
		log.Printf("Error updating borrow record for return: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update borrow record"})
	}

	var book models.Book
	if err := db.DB.First(&book, borrow.BookID).Error; err != nil {
		log.Printf("Error finding book to update availability after return: %v", err)
	} else {
		if err := db.DB.Model(&book).Update("available", true).Error; err != nil {
			log.Printf("Error updating book availability after return: %v", err)
		}
	}

	var user models.User
	if err := db.DB.First(&user, borrow.UserID).Error; err != nil {
		log.Printf("Error finding user to update penalty after return: %v", err)
	} else {
		user.Penalty += fineAmount 
		if err := db.DB.Save(&user).Error; err != nil {
			log.Printf("Error updating user penalty: %v", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Book returned successfully",
		"borrow_id":   borrow.ID,
		"book_id":     borrow.BookID,
		"user_id":     borrow.UserID,
		"fine_incurred": fineAmount,
		"is_overdue":  fineAmount > 0,
	})
}