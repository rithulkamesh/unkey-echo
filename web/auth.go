package web

import (
	"net/http"
	"time"

	"github.com/kordlab/marketplace/data"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func registerAuthRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	auth.POST("/login", handleLogin)
	auth.POST("/register", handleRegister)
	auth.POST("/logout", handleLogout)
}

func handleLogin(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// TODO: Implement database lookup and password verification
	// For now returning mock response
	user := &data.User{
		ID:        primitive.NewObjectID(),
		Email:     req.Email,
		Role:      data.RoleUser,
		Status:    data.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return c.JSON(http.StatusOK, user)
}

func handleRegister(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process password"})
	}

	// Create new user
	user := &data.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         data.RoleUser,
		Status:       data.UserStatusActive,
		Profile: data.UserProfile{
			DisplayName: req.Username,
		},
		Credits: data.Credits{
			Balance: 0,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// TODO: Implement database insertion

	// For now returning mock response
	return c.JSON(http.StatusCreated, user)
}

func handleLogout(c echo.Context) error {
	// TODO: Implement session/token invalidation
	return c.NoContent(http.StatusOK)
}
