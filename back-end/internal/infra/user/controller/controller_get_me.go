package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/exception"
	"github.com/laurencefluciano/content-api/internal/infra"
	"github.com/laurencefluciano/content-api/internal/infra/middleware"
	"github.com/laurencefluciano/content-api/internal/infra/user/dto"
)

func (h *UserController) getMe(c *gin.Context) {
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

	userEntity, err := h.repo.FindByAuthId(userClaims.AuthID)

	if err != nil {
		if errors.Is(err, exception.ErrRepoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		}

		h.handleAppError(c, err)
		return
	}

	response := dto.ProfileResponse{
		ID:        userEntity.ID().String(),
		Name:      infra.StringPtr(userEntity.Name()),
		Roles:     user.ToStringRoles(userEntity.Roles()),
		AvatarUrl: infra.StringPtr(userEntity.AvatarUrl()),
	}

	if userEntity.HasRole(user.RoleProducer) {
		response.Producer = &dto.ProducerProfileResponse{
			Bio:      infra.StringPtr(userEntity.Bio()),
			Websites: userEntity.Websites(),
			Social:   userEntity.Social(),
		}
	}

	c.JSON(http.StatusOK, response)
}
