// router/router.go
package router

import (
	propertyService "github.com/PauloGuillen/gostosobookings/internal/property/service"
	userService "github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/PauloGuillen/gostosobookings/pkg/controller"
	"github.com/PauloGuillen/gostosobookings/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	usrService userService.UserService, authService userService.AuthService, propService propertyService.PropertyService,
) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		userController := controller.NewUserController(usrService)
		authController := controller.NewAuthController(authService)
		propertyController := controller.NewPropertyController(propService)

		v1.POST("/login", authController.Login)
		v1.POST("/logout", authController.Logout)
		v1.POST("/users", userController.CreateUser)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))

		protected.PUT("/users/:id", userController.UpdateUser)

		protected.POST("/property", propertyController.CreateProperty)

		// Endpoint para Promover customer para business_admin
		// PUT /users/{id}/role
		// func UpdateUserRole(c *gin.Context) {

		// Endpoint para Promover customer para business_manager:
		// Restrito a business_admin relacionado ao negócio.
		// PUT /businesses/{business_id}/users/{id}/role
		// func UpdateBusinessUserRole(c *gin.Context) {

	}

	return r
}
