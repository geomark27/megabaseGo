package handlers

import (
	"net/http"
	"strconv"

	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/app/services"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *services.RoleService
}

// NewRoleHandler crea una nueva instancia del handler de roles
func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: services.NewRoleService(),
	}
}

// CreateRole maneja la creación de roles
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		if err.Error() == "role with this name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create role",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"role":    role,
	})
}

// GetRoles maneja la obtención de todos los roles
func (h *RoleHandler) GetRoles(c *gin.Context) {
	includeInactive := c.Query("include_inactive") == "true"

	roles, err := h.roleService.GetRoles(includeInactive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch roles",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
		"count": len(roles),
	})
}

// GetRole maneja la obtención de un rol por ID
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	role, err := h.roleService.GetRoleByID(uint(roleID))
	if err != nil {
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch role",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"role": role,
	})
}

// UpdateRole maneja la actualización de roles
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	role, err := h.roleService.UpdateRole(uint(roleID), &req)
	if err != nil {
		switch err.Error() {
		case "role not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "role with this name already exists":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update role",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
		"role":    role,
	})
}

// DeleteRole maneja la eliminación de roles
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	err = h.roleService.DeleteRole(uint(roleID))
	if err != nil {
		switch err.Error() {
		case "role not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "cannot delete role: it is assigned to users":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to delete role",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role deleted successfully",
	})
}
