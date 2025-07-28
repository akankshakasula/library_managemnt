package models

import (
	"gorm.io/gorm"
)

// Define roles as constants for clarity and type safety
const (
	RoleLibrarian = "librarian"
	RoleStudent   = "student"
	RoleGeneral   = "general"
)

type User struct {
	gorm.Model
	Name     string  `json:"name"`
	Email    string  `json:"email" gorm:"unique"`
	Password string  `json:"-"` // Don't return password in JSON
	Role     string  `json:"role"`
	Penalty  float64 `json:"penalty" gorm:"default:0.0"` // Total accumulated fines
	Blocked  bool    `json:"blocked" gorm:"default:false"`
}

// IsValidRole checks if the given role string is one of the predefined valid roles.
func IsValidRole(role string) bool {
	switch role {
	case RoleLibrarian, RoleStudent, RoleGeneral:
		return true
	default:
		return false
	}
}