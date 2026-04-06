package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("", h.register)
		userRoutes.POST("webhook/confirm", h.confirmWebhook)
		userRoutes.GET("/me", h.getMe)
		userRoutes.PATCH("/me", h.updateMe)
		userRoutes.POST("/become-producer", h.becomeProducer)
	}
}

func (h *Handler) register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) getMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) updateMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) confirmWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) becomeProducer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
