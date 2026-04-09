package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func CreateHandler(r *gin.RouterGroup) *Handler {
	handler := &Handler{}
	handler.InitRoutes(r)
	return handler
}

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("", AuthMiddleware(), h.register)
		userRoutes.POST("webhook/confirm", h.confirmWebhook)
		userRoutes.GET("/me", h.getMe)
		userRoutes.PATCH("/me", h.updateMe)
		userRoutes.POST("/become-producer", h.becomeProducer)
	}
}

func (h *Handler) register(c *gin.Context) {
	var dto CreateUserRequest

	err := c.ShouldBindJSON(&dto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, ok := c.Get("auth_id")

	if !ok {
		log.Println("[ ERRO ]: 'auth_id' não encontrado no contexto")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	authId, ok := val.(string)
	if !ok {
		log.Println("[ ERRO ]: Erro ao converter o tipo para string")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro"})
		return
	}

	userId, err := h.repo.CreateUser(&SaveUserParams{
		name:   dto.Name,
		authId: authId,
		status: "pending",
	})

	if err != nil {
		log.Printf("[ ERRO ]: Falha ao persistir usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar no banco de dados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user register",
		"status":  "created",
		"id":      userId,
	})
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
