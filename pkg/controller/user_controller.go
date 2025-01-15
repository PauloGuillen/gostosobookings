package controller

import (
	"fmt"
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/helpers"
	"github.com/PauloGuillen/gostosobookings/internal/user/dto"
	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/gin-gonic/gin"
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
		fmt.Println("ShouldBindJSON - error", err)
		if handled := helpers.HandleValidationError(ctx, err); handled {
			fmt.Println("ShouldBindJSON - handled", handled)
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
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
