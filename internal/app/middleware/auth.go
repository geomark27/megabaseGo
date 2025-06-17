package middleware

import (
	"net/http"
	"strings"

	"megabaseGo/internal/app/services"
	"megabaseGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware maneja la autenticación JWT
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware crea una nueva instancia del middleware de autenticación
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		authService: services.NewAuthService(),
	}
}

// RequireAuth middleware que requiere autenticación
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.UserName)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID)
		c.Set("role_name", claims.RoleName)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole middleware que requiere un rol específico
func (m *AuthMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primero verificar autenticación
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Verificar rol
		userRoleName, exists := c.Get("role_name")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Role information not found",
			})
			c.Abort()
			return
		}

		if userRoleName != roleName {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware que requiere uno de varios roles
func (m *AuthMiddleware) RequireAnyRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primero verificar autenticación
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Verificar si el usuario tiene alguno de los roles permitidos
		userRoleName, exists := c.Get("role_name")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Role information not found",
			})
			c.Abort()
			return
		}

		roleMatches := false
		for _, roleName := range roleNames {
			if userRoleName == roleName {
				roleMatches = true
				break
			}
		}

		if !roleMatches {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware que permite autenticación opcional
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No hay token, continuar sin autenticación
			c.Next()
			return
		}

		// Hay token, intentar validarlo
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			if claims, err := m.authService.ValidateToken(tokenString); err == nil {
				// Token válido, guardar información en contexto
				c.Set("user_id", claims.UserID)
				c.Set("user_name", claims.UserName)
				c.Set("email", claims.Email)
				c.Set("role_id", claims.RoleID)
				c.Set("role_name", claims.RoleName)
				c.Set("claims", claims)
				c.Set("authenticated", true)
			}
		}

		c.Next()
	}
}

// GetCurrentUserID obtiene el ID del usuario actual del contexto
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentUserClaims obtiene los claims del usuario actual
func GetCurrentUserClaims(c *gin.Context) (*utils.JWTClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	return claims.(*utils.JWTClaims), true
}

// IsAuthenticated verifica si el usuario está autenticado
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}