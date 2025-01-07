package router

import (
	"github.com/PauloGuillen/gostosobookings/pkg/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("/users", controller.CreateUser)
		// v1.GET("/users/:id", controller.GetUserByID)
		// v1.PUT("/users/:id", controller.UpdateUser)
		// v1.DELETE("/users/:id", controller.DeleteUser)
		// v1.POST("/login", controller.Login)
		// v1.POST("/logout", controller.Logout)
	}

	return r
}
