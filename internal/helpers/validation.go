package helpers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleValidationError processes validation errors and sends a structured response.
func HandleValidationError(ctx *gin.Context, err error) bool {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string]string)

		// Iterate over validation errors and map field-specific messages
		for _, e := range validationErrs {
			switch e.Tag() {
			case "required":
				fieldErrors[e.Field()] = fmt.Sprintf("%s is required", e.Field())
			case "email":
				fieldErrors[e.Field()] = "Invalid email format"
			case "min":
				fieldErrors[e.Field()] = fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
			default:
				fieldErrors[e.Field()] = fmt.Sprintf("Invalid value for %s", e.Field())
			}
		}

		// Return 422 with validation errors
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Validation error",
			"errors":  fieldErrors,
		})
		return true
	}

	// Return false if the error is not a validation error
	return false
}
