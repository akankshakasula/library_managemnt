package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string `json:"title"`
	Author      string `json:"author"`
	Number      string `json:"number" gorm:"unique"` 
	Genre       string `json:"genre"`
	DonatedByID uint   `json:"donated_by_id"`        
	DonatedBy   User   `json:"-" gorm:"foreignKey:DonatedByID"` 
	Available   bool   `json:"available" gorm:"default:true"` 
}