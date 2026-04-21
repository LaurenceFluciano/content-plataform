package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/app/user/command"
	"github.com/laurencefluciano/content-api/internal/exception"
	"github.com/laurencefluciano/content-api/internal/infra/middleware"
)

func (h *UserController) register(c *gin.Context) {
	userClaims, ok := middleware.GetUserAuth(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
	}

	if !userClaims.EmailVerified {
		c.JSON(http.StatusNotFound, &exception.AppError{
			Code:    exception.ValidationFailedCode,
			Message: "Você precisa verificar seu e-mail para se tornar um produtor.",
		})
		return
	}

	err := h.registerUser.Execute(command.RegisterUserCommand{
		AuthId: userClaims.AuthID,
	})

	if err != nil {
		h.handleAppError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created.",
	})
}
