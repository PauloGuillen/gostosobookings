package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Aqui você adiciona lógica para validar o email e senha,
	// e gerar o token JWT se estiver tudo correto.
	// Exemplo básico de resposta:
	token := "generated_jwt_token" // Substituir pela lógica real de geração de token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	// Implementar a lógica de logout, como invalidar o token se for aplicável
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
