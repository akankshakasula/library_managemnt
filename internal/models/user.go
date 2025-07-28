package models

import (
	"gorm.io/gorm"
)

const (
	RoleLibrarian = "librarian"
	RoleStudent   = "student"
	RoleGeneral   = "general"
)

type User struct {
	gorm.Model
	Name     string  `json:"name"`
	Email    string  `json:"email" gorm:"unique"`
	Password string  `json:"-"` 
	Role     string  `json:"role"`
	Penalty  float64 `json:"penalty" gorm:"default:0.0"` 
	Blocked  bool    `json:"blocked" gorm:"default:false"`
}

func IsValidRole(role string) bool {
	switch role {
	case RoleLibrarian, RoleStudent, RoleGeneral:
		return true
	default:
		return false
	}
}