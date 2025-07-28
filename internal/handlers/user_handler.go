package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // CHANGE THIS LINE BACK TO v5
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"library-management/internal/db"
	"library-management/internal/models"
)

// SignUpRequest struct for user registration
type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // e.g., "librarian", "student", "general"
}

// SignInRequest struct for user login
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// --- JWT Claims (Payload) ---
// Note: UserID is uint as per your model, not string.
type Claims struct {
	UserID               uint   `json:"user_id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // This usually provides GetAudience, GetExpiresAt etc.
}

// You DO NOT need to explicitly implement GetAudience, GetExpiresAt etc. here
// if you are correctly using jwt.RegisteredClaims from github.com/golang-jwt/jwt/v5.
// The embedding should handle it automatically.
// If after all steps, it still complains, then we can re-add them as a workaround.

// SignUp handles user registration
func SignUp(c *fiber.Ctx) error {
	req := new(SignUpRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, Email, Password, and Role are required"})
	}

	// Validate role
	if !models.IsValidRole(req.Role) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role. Must be 'librarian', 'student', or 'general'"})
	}

	// Check if user with this email already exists
	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User with this email already exists"})
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Database error checking for existing user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Hash the password
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
		// Penalty and Blocked will default to 0.0 and false
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

// SignIn handles user login and generates a JWT
func SignIn(c *fiber.Ctx) error {
	req := new(SignInRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Basic validation
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

	// Check if user is blocked
	if user.Blocked {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Your account is blocked. Please contact the librarian."})
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// --- Generate JWT Token ---
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
		"token":      tokenString, // Return the JWT token
		"user_id":    user.ID,
		"user_name":  user.Name,
		"user_email": user.Email,
		"user_role":  user.Role,
	})
}
