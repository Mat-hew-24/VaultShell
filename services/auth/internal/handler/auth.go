package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"syncpad/services/auth/internal/service"
	"syncpad/services/auth/pkg/jwt"
)

// AuthHandler handles auth-related HTTP requests
type AuthHandler struct {
	auth *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register endpoint
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Call AuthService to register user
	err := h.auth.Register(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		} else if errors.Is(err, service.ErrInvalidEmail) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		} else if errors.Is(err, service.ErrPasswordTooShort) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
		} else if errors.Is(err, service.ErrEmptyField) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Success
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// Login endpoint
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Call AuthService to login user
	user, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
