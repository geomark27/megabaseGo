package routes

import (
	"megabaseGo/internal/app/handlers"
	"megabaseGo/internal/app/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup configura todas las rutas de la aplicación
func Setup() *gin.Engine {
	// Crear router con configuración por defecto
	router := gin.Default()

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	router.Use(cors.New(config))

	// Middleware de logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Ruta de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "MegabaseGo API is running",
			"version": "1.0.0",
		})
	})

	// Inicializar handlers y middleware
	userHandler := handlers.NewUserHandler()
	roleHandler := handlers.NewRoleHandler()
	authHandler := handlers.NewAuthHandler()
	authMiddleware := middleware.NewAuthMiddleware()

	// Grupo de rutas API v1
	v1 := router.Group("/api/v1")
	{
		// Rutas públicas de autenticación
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)                    // POST /api/v1/auth/login
			auth.POST("/register", authHandler.Register)              // POST /api/v1/auth/register
			auth.POST("/refresh", authHandler.RefreshToken)           // POST /api/v1/auth/refresh
			auth.POST("/logout", authHandler.Logout)                  // POST /api/v1/auth/logout
		}

		// Rutas protegidas (requieren autenticación)
		protected := v1.Group("/")
		protected.Use(authMiddleware.RequireAuth())
		{
			// Profile endpoints
			protected.GET("/profile", authHandler.GetProfile)           // GET /api/v1/profile
			protected.POST("/change-password", authHandler.ChangePassword) // POST /api/v1/change-password
			protected.GET("/check-auth", authHandler.CheckAuth)         // GET /api/v1/check-auth

			// Rutas para roles (requiere autenticación)
			roles := protected.Group("/roles")
			{
				roles.POST("", roleHandler.CreateRole)           // POST /api/v1/roles
				roles.GET("", roleHandler.GetRoles)              // GET /api/v1/roles
				roles.GET("/:id", roleHandler.GetRole)           // GET /api/v1/roles/:id
				roles.PUT("/:id", roleHandler.UpdateRole)        // PUT /api/v1/roles/:id
				roles.DELETE("/:id", roleHandler.DeleteRole)     // DELETE /api/v1/roles/:id
			}

			// Rutas para usuarios (requiere autenticación)
			users := protected.Group("/users")
			{
				users.POST("", userHandler.CreateUser)           // POST /api/v1/users
				users.GET("", userHandler.GetUsers)              // GET /api/v1/users
				users.GET("/:id", userHandler.GetUser)           // GET /api/v1/users/:id
				users.PUT("/:id", userHandler.UpdateUser)        // PUT /api/v1/users/:id
				users.DELETE("/:id", userHandler.DeleteUser)     // DELETE /api/v1/users/:id
			}
		}

		// Ruta de información de la API (pública)
		v1.GET("/info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"api_version": "1.0.0",
				"project":     "MegabaseGo",
				"architecture": gin.H{
					"pattern": "Handler + Service (Simplified) + JWT Auth",
					"description": "Clean architecture with JWT authentication for rapid development",
				},
				"endpoints": gin.H{
					"auth": gin.H{
						"login":     "POST /api/v1/auth/login",
						"register":  "POST /api/v1/auth/register",
						"refresh":   "POST /api/v1/auth/refresh",
						"logout":    "POST /api/v1/auth/logout",
						"profile":   "GET /api/v1/profile (protected)",
						"check":     "GET /api/v1/check-auth (protected)",
						"password":  "POST /api/v1/change-password (protected)",
					},
					"roles": gin.H{
						"create": "POST /api/v1/roles (protected)",
						"list":   "GET /api/v1/roles (protected)",
						"get":    "GET /api/v1/roles/:id (protected)",
						"update": "PUT /api/v1/roles/:id (protected)",
						"delete": "DELETE /api/v1/roles/:id (protected)",
					},
					"users": gin.H{
						"create": "POST /api/v1/users (protected)",
						"list":   "GET /api/v1/users (protected)",
						"get":    "GET /api/v1/users/:id (protected)",
						"update": "PUT /api/v1/users/:id (protected)",
						"delete": "DELETE /api/v1/users/:id (protected)",
					},
				},
				"authentication": gin.H{
					"type":   "JWT Bearer Token",
					"header": "Authorization: Bearer <token>",
					"note":   "Include access token in Authorization header for protected routes",
				},
				"query_params": gin.H{
					"roles": gin.H{
						"include_inactive": "bool - Include inactive roles",
					},
					"users": gin.H{
						"include_inactive": "bool - Include inactive users",
						"role_id":          "int - Filter by role ID",
					},
				},
			})
		})
	}

	return router
}