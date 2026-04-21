package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/app/user/command"
	"github.com/laurencefluciano/content-api/internal/infra/middleware"
)

func (h *UserController) becomeProducer(c *gin.Context) {
	userClaims, ok := middleware.GetUserAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	err := h.becameProducer.Execute(command.BecomeProducerCommand{
		AuthId: userClaims.AuthID,
	})

	if err != nil {
		h.handleAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parabéns! Você agora é um produtor."})
}
