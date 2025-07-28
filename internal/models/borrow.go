package models

import (
	"time"

	"gorm.io/gorm"
)

type Borrow struct {
	gorm.Model
	BookID       uint      `json:"book_id"`
	Book         Book      `json:"-" gorm:"foreignKey:BookID"` // GORM association for the borrowed book
	UserID       uint      `json:"user_id"`
	User         User      `json:"-" gorm:"foreignKey:UserID"` // GORM association for the user who borrowed
	BorrowDate   time.Time `json:"borrow_date"`
	ReturnDate   *time.Time `json:"return_date"` // Pointer because it can be null until returned
	DueDate      time.Time `json:"due_date"`
	Returned     bool      `json:"returned" gorm:"default:false"`
	FineAmount   float64   `json:"fine_amount" gorm:"default:0.0"`
	FinePaid     bool      `json:"fine_paid" gorm:"default:false"`
}