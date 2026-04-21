package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/app/user/command"
	"github.com/laurencefluciano/content-api/internal/exception"
	"github.com/laurencefluciano/content-api/internal/infra/middleware"
	"github.com/laurencefluciano/content-api/internal/infra/user/dto"
)

func (h *UserController) updateMe(c *gin.Context) {
	userClaims, ok := middleware.GetUserAuth(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	if !userClaims.EmailVerified {
		c.JSON(http.StatusNotFound, &exception.AppError{
			Code:    exception.ValidationFailedCode,
			Message: "Você precisa verificar seu e-mail para se tornar um produtor.",
		})
		return
	}

	var dto dto.UpdateUserRequest

	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	err := h.updateProfile.Execute(userClaims.AuthID, command.UpdateUserCommand{
		Name:      dto.Name,
		AvatarUrl: dto.AvatarUrl,
		Bio:       dto.Bio,
		Websites:  dto.Websites,
		Social:    dto.Social,
	})

	if err != nil {
		h.handleAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil de usuário atualizado."})
}
