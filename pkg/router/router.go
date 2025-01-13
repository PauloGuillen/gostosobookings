// router/router.go
package router

import (
	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/PauloGuillen/gostosobookings/pkg/controller"
	"github.com/PauloGuillen/gostosobookings/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userService service.UserService, authService service.AuthService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		userController := controller.NewUserController(userService)
		authController := controller.NewAuthController(authService)

		v1.POST("/login", authController.Login)
		v1.POST("/logout", authController.Logout)
		v1.POST("/users", userController.CreateUser)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))

		protected.PUT("/users/:id", userController.UpdateUser)
	}

	return r
}
