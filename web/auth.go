package web

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kordlab/marketplace/data"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
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

type AuthHandler struct {
	mongo *data.MongoDB
	redis *data.RedisDB
}

func NewAuthHandler(mongo *data.MongoDB, redis *data.RedisDB) *AuthHandler {
	return &AuthHandler{
		mongo: mongo,
		redis: redis,
	}
}

func registerAuthRoutes(e *echo.Echo, h *AuthHandler) {
	auth := e.Group("/auth")
	auth.POST("/login", h.handleLogin)
	auth.POST("/register", h.handleRegister)
	auth.POST("/logout", h.handleLogout)
}

func (h *AuthHandler) handleLogin(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Find user by email
	var user data.User
	err := h.mongo.Users().FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Check user status
	if user.Status != data.UserStatusActive {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Account is not active"})
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID.Hex()
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(h.mongo.GetJWTSecret()))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	// Store session in Redis
	err = h.redis.StoreSession(user.ID.Hex(), tokenString, time.Hour*24)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": tokenString,
		"user":  user,
	})
}

func (h *AuthHandler) handleRegister(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if email already exists
	exists, err := h.mongo.Users().CountDocuments(context.Background(), bson.M{"email": req.Email})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if exists > 0 {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email already registered"})
	}

	// Check if username already exists
	exists, err = h.mongo.Users().CountDocuments(context.Background(), bson.M{"username": req.Username})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if exists > 0 {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Username already taken"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process password"})
	}

	now := time.Now()
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
			Balance:      0,
			Transactions: []data.CreditTransaction{},
		},
		Notifications: []data.Notification{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Insert user into database
	_, err = h.mongo.Users().InsertOne(context.Background(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	// Cache the new user
	err = h.redis.CacheUser(user)
	if err != nil {
		// Log the error but don't fail the request
		c.Logger().Error("Failed to cache user:", err)
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) handleLogout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.NoContent(http.StatusOK)
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Blacklist the token
	err := h.redis.BlacklistToken(token, time.Hour*24)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to logout"})
	}

	return c.NoContent(http.StatusOK)
}
