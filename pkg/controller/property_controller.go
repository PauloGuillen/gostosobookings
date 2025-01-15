package controller

import (
	"github.com/PauloGuillen/gostosobookings/internal/property/service"
	"github.com/gin-gonic/gin"
)

type PropertyController struct {
	propertyService service.PropertyService
}

func NewPropertyController(propertyService service.PropertyService) *PropertyController {
	return &PropertyController{propertyService: propertyService}
}

// CreateProperty handles the creation of a new property.
func (c *PropertyController) CreateProperty(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	role, _ := ctx.Get("role")

	ctx.JSON(200, gin.H{
		"message": "CreateProperty",
		"user_id": userID,
		"role":    role,
	})

}
