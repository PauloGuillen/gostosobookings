package controller

import (
	"fmt"
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/helpers"
	"github.com/PauloGuillen/gostosobookings/internal/property/dto"
	"github.com/PauloGuillen/gostosobookings/internal/property/service"
	"github.com/gin-gonic/gin"
)

type PropertyController struct {
	propService service.PropertyService
}

func NewPropertyController(propertyService service.PropertyService) *PropertyController {
	return &PropertyController{propService: propertyService}
}

// CreateProperty handles the creation of a new property.
func (c *PropertyController) CreateProperty(ctx *gin.Context) {
	var propertyRequest dto.CreatePropertyRequest

	// Bind JSON to propertyRequest and check for errors
	if err := ctx.ShouldBindJSON(&propertyRequest); err != nil {
		fmt.Println("ShouldBindJSON - error", err)
		if handled := helpers.HandleValidationError(ctx, err); handled {
			fmt.Println("ShouldBindJSON - handled", handled)
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	userID, _ := ctx.Get("user_id")
	role, _ := ctx.Get("role")

	// Call the service to create the property
	// err := c.propService.CreateProperty(ctx.Request.Context(), userID, role)

	ctx.JSON(200, gin.H{
		"message": "CreateProperty",
		"user_id": userID,
		"role":    role,
	})

}
