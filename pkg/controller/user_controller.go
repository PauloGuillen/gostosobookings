package controller

import (
	"fmt"
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/dto"
	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userService service.UserService
}

// NewUserController creates a new instance of UserController.
func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// CreateUser handles the creation of a new user.
func (c *UserController) CreateUser(ctx *gin.Context) {
	var userRequest dto.CreateUserRequest
	// Bind JSON to userRequest and check for errors
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
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
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Validation error",
				"errors":  fieldErrors,
			})
			return
		}
	}

	// Proceed with user creation if validation passes
	user, err := c.userService.CreateUser(ctx.Request.Context(), userRequest)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

// UpdateUser handles the update of an existing user.
func (c *UserController) UpdateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "UpdateUser"})
}
