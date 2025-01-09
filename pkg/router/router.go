// router/router.go
package router

import (
	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/PauloGuillen/gostosobookings/pkg/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userService service.UserService, authService service.AuthService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		userController := controller.NewUserController(userService)
		authController := controller.NewAuthController(authService)

		v1.POST("/users", userController.CreateUser)
		v1.POST("/login", authController.Login)
	}

	return r
}
