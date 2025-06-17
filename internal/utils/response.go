package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError maneja errores con status codes automáticos
func HandleError(c *gin.Context, err error) {
	if apiErr, ok := IsAPIError(err); ok {
		c.JSON(apiErr.GetStatusCode(), gin.H{"error": apiErr.Message})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

// HandleValidationError maneja errores de validación
func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Invalid request data",
		"details": err.Error(),
	})
}

// HandleSuccess maneja respuestas exitosas
func HandleSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := gin.H{"message": message}
	
	if data != nil {
		if dataMap, ok := data.(gin.H); ok {
			for key, value := range dataMap {
				response[key] = value
			}
		} else {
			response["data"] = data
		}
	}
	
	c.JSON(statusCode, response)
}

// HandleData maneja respuesta con datos sin mensaje
func HandleData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}