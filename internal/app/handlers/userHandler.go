package handlers

import (
	"net/http"
	"strconv"

	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/app/services"
	"megabaseGo/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler crea una nueva instancia del handler de usuarios
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// CreateUser maneja la creación de usuarios
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		utils.HandleError(c, err) // ✅ Una sola línea!
		return
	}

	utils.HandleSuccess(c, http.StatusCreated, "User created successfully", gin.H{"user": user})
}

// GetUsers maneja la obtención de todos los usuarios
func (h *UserHandler) GetUsers(c *gin.Context) {
	includeInactive := c.Query("include_inactive") == "true"
	
	var roleID *uint
	if roleIDStr := c.Query("role_id"); roleIDStr != "" {
		if id, err := strconv.ParseUint(roleIDStr, 10, 32); err == nil {
			roleIDUint := uint(id)
			roleID = &roleIDUint
		}
	}

	users, err := h.userService.GetUsers(includeInactive, roleID)
	if err != nil {
		utils.HandleError(c, err) // ✅ Una sola línea!
		return
	}

	utils.HandleData(c, http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// GetUser maneja la obtención de un usuario por ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.NewBadRequestError("Invalid user ID"))
		return
	}

	user, err := h.userService.GetUserByID(uint(userID))
	if err != nil {
		utils.HandleError(c, err) // ✅ Una sola línea!
		return
	}

	utils.HandleData(c, http.StatusOK, gin.H{"user": user})
}

// UpdateUser maneja la actualización de usuarios
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.NewBadRequestError("Invalid user ID"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	user, err := h.userService.UpdateUser(uint(userID), &req)
	if err != nil {
		utils.HandleError(c, err) // ✅ Una sola línea!
		return
	}

	utils.HandleSuccess(c, http.StatusOK, "User updated successfully", gin.H{"user": user})
}

// DeleteUser maneja la eliminación de usuarios
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleError(c, utils.NewBadRequestError("Invalid user ID"))
		return
	}

	err = h.userService.DeleteUser(uint(userID))
	if err != nil {
		utils.HandleError(c, err) // ✅ Una sola línea!
		return
	}

	utils.HandleSuccess(c, http.StatusOK, "User deleted successfully", nil)
}