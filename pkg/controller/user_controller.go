package controller

import (
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/dto"
	"github.com/PauloGuillen/gostosobookings/internal/user/repository"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var userRequest dto.CreateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
	}

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
