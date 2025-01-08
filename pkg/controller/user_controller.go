package controller

import (
	"fmt"
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/dto"
	"github.com/PauloGuillen/gostosobookings/internal/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateUser(c *gin.Context) {
	var userRequest dto.CreateUserRequest
	// Bind JSON to userRequest and check for errors
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		// Map for field-specific errors
		fieldErrors := make(map[string]string)

		// Iterate over validation errors
		for _, e := range err.(validator.ValidationErrors) {
			// Custom error messages for each field
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

		// Return 422 status with field-specific errors
		if len(fieldErrors) > 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Validation error",
				"errors":  fieldErrors,
			})
			return
		}
	}

	// Proceed with user creation if validation passes
	user, err := repository.CreateUser(userRequest)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// handleError is a helper function to handle different types of errors and send the appropriate HTTP response
func handleError(c *gin.Context, err error) {
	var statusCode int
	var errorMessage string

	switch err {
	case errors.ErrEmailAlreadyExists:
		statusCode = http.StatusConflict
		errorMessage = err.Error()
	case errors.ErrPasswordHashing, errors.ErrSonyflakeInit, errors.ErrSonyflakeNextID, errors.ErrDatabase:
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = "Unknown error"
	}

	c.JSON(statusCode, gin.H{"error": errorMessage})
}
