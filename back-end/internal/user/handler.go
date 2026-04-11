package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	repo *Repository
}

func CreateHandler(r *gin.RouterGroup, repo *Repository) *Handler {
	handler := &Handler{repo: repo}
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

	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var name string

	if dto.Name != nil {
		if len(*dto.Name) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nome deve ter no minimo 2 caracteres."})
			return
		}
		name = *dto.Name
	}

	userClaims, ok := GetUserAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	status := Pending
	if userClaims.EmailVerified {
		status = Active
	}

	user, err := h.repo.GetUserByAuthId(userClaims.AuthID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userId, createErr := h.repo.CreateUser(&SaveUserParams{
				name:   name,
				authId: userClaims.AuthID,
				status: status,
			})

			if createErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno"})
				return
			}

			err := h.repo.AddRole(uuid.MustParse(userId), RoleReader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atribuir permissão"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"id": userId, "status": "created"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro de conexão com o banco"})
		return
	}

	if !user.HasRole("reader") {
		err := h.repo.AddRole(user.ID, RoleReader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atribuir permissão"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     user.ID,
		"status": status,
	})
}

func (h *Handler) getMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})

	userClaims, ok := GetUserAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	profile, err := h.repo.GetProfileByAuthId(userClaims.AuthID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	response := gin.H{
		"id":   profile.User.ID,
		"name": profile.User.Name,
	}

	if profile.UserProfile != nil {
		response["avatar_url"] = profile.UserProfile.AvatarUrl
	}

	if profile.ProducerProfile != nil {
		response["producer"] = gin.H{
			"bio":      profile.ProducerProfile.Bio,
			"websites": profile.ProducerProfile.Websites,
			"social":   profile.ProducerProfile.Social,
		}
	}

	c.JSON(http.StatusOK, response)
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
