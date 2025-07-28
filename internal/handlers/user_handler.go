package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" 
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"library-management/internal/db"
	"library-management/internal/models"
)

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` 
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID               uint   `json:"user_id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        
}


func SignUp(c *fiber.Ctx) error {
	req := new(SignUpRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, Email, Password, and Role are required"})
	}

	if !models.IsValidRole(req.Role) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role. Must be 'librarian', 'student', or 'general'"})
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User with this email already exists"})
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Database error checking for existing user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func SignIn(c *fiber.Ctx) error {
	req := new(SignInRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		log.Printf("Database error finding user for login: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if user.Blocked {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Your account is blocked. Please contact the librarian."})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(db.JWTSecret)) // Use the secret from db package
	if err != nil {
		log.Printf("Error signing JWT token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Login successful",
		"token":      tokenString, 
		"user_id":    user.ID,
		"user_name":  user.Name,
		"user_email": user.Email,
		"user_role":  user.Role,
	})
}
