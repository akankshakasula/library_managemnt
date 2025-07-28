package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string `json:"title"`
	Author      string `json:"author"`
	Number      string `json:"number" gorm:"unique"` // Unique code/number for the book
	Genre       string `json:"genre"`
	DonatedByID uint   `json:"donated_by_id"`        // Foreign key to User, using `uint` for GORM ID
	DonatedBy   User   `json:"-" gorm:"foreignKey:DonatedByID"` // GORM association, exclude from JSON
	Available   bool   `json:"available" gorm:"default:true"` // true if available for borrowing
}