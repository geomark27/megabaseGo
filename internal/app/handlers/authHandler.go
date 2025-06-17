package handlers

import (
	"net/http"

	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/app/middleware"
	"megabaseGo/internal/app/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler crea una nueva instancia del handler de autenticación
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Login maneja el inicio de sesión
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		switch err.Error() {
		case "invalid credentials":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		case "user account is disabled":
			c.JSON(http.StatusForbidden, gin.H{"error": "User account is disabled"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Login failed",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    authResponse,
	})
}

// Register maneja el registro de nuevos usuarios
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		switch err.Error() {
		case "role not found":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
		case "username already exists", "email already exists":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Registration failed",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"data":    authResponse,
	})
}

// RefreshToken maneja la renovación de tokens
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	authResponse, err := h.authService.RefreshToken(&req)
	if err != nil {
		switch err.Error() {
		case "invalid refresh token":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case "user account is disabled":
			c.JSON(http.StatusForbidden, gin.H{"error": "User account is disabled"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Token refresh failed",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data":    authResponse,
	})
}

// Logout maneja el cierre de sesión (invalidar token del lado del cliente)
func (h *AuthHandler) Logout(c *gin.Context) {
	// En este caso simple, el logout es manejado por el cliente
	// removiendo el token. En implementaciones más avanzadas,
	// podrías mantener una blacklist de tokens invalidados.
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// GetProfile obtiene el perfil del usuario actual
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch user profile",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// ChangePassword maneja el cambio de contraseña
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	err := h.authService.ChangePassword(userID, &req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case "current password is incorrect":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to change password",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// CheckAuth verifica si el usuario está autenticado
func (h *AuthHandler) CheckAuth(c *gin.Context) {
	claims, exists := middleware.GetCurrentUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user": gin.H{
			"id":        claims.UserID,
			"user_name": claims.UserName,
			"email":     claims.Email,
			"role_id":   claims.RoleID,
			"role_name": claims.RoleName,
		},
	})
}