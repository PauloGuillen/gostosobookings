package router

import (
	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/PauloGuillen/gostosobookings/pkg/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userService service.UserService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		// Criar uma instância do controller com o userService
		userController := controller.NewUserController(userService)

		// Usar a função CreateUser do controller
		v1.POST("/users", userController.CreateUser)
		// Adicionar outros endpoints...
	}

	return r
}
