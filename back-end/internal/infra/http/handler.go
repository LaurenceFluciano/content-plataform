package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/app/exception"
	"github.com/laurencefluciano/content-api/internal/app/user/usecase"
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/infra"
	"github.com/laurencefluciano/content-api/internal/infra/middleware"
)

type Handler struct {
	registerUser   usecase.RegisterUserUseCase
	updateProfile  usecase.UpdateProfileUseCase
	becameProducer usecase.BecomeProducerUserCase
	repo           user.Repository
}

func CreateHandler(r *gin.RouterGroup, repo user.Repository) *Handler {
	handler := &Handler{
		registerUser: usecase.RegisterUserUseCase{
			Repo: repo,
		},
		updateProfile: usecase.UpdateProfileUseCase{
			Repo: repo,
		},
		becameProducer: usecase.BecomeProducerUserCase{
			Repo: repo,
		},
		repo: repo,
	}
	handler.InitRoutes(r)
	return handler
}

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/", middleware.AuthMiddleware(), h.register)
		userRoutes.GET("/me/profile", middleware.AuthMiddleware(), h.getMe)
		userRoutes.PATCH("/me/profile", middleware.AuthMiddleware(), h.updateMe)
		userRoutes.POST("/become-producer", middleware.AuthMiddleware(), h.becomeProducer)
	}
}

func (h *Handler) register(c *gin.Context) {
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

	err := h.registerUser.Execute(userClaims.AuthID)

	if err != nil {
		h.handleAppError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created.",
	})
}

func (h *Handler) getMe(c *gin.Context) {
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
		if errors.Is(err, exception.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		}

		h.handleAppError(c, err)
		return
	}

	response := ProfileResponse{
		ID:        userEntity.ID().String(),
		Name:      infra.StringPtr(userEntity.Name()),
		Roles:     user.ToStringRoles(userEntity.Roles()),
		AvatarUrl: infra.StringPtr(userEntity.Profile().AvatarUrl()),
	}

	if userEntity.HasRole(user.RoleProducer) {
		response.Producer = &ProducerDTO{
			Bio:      infra.StringPtr(userEntity.Profile().Bio()),
			Websites: userEntity.Profile().Websites(),
			Social:   userEntity.Profile().Social(),
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) updateMe(c *gin.Context) {
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

	var dto UpdateUserRequest

	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	err := h.updateProfile.Execute(userClaims.AuthID, usecase.UpdateUserParams{
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

func (h *Handler) becomeProducer(c *gin.Context) {
	userClaims, ok := middleware.GetUserAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	err := h.becameProducer.Execute(userClaims.AuthID)

	if err != nil {
		h.handleAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parabéns! Você agora é um produtor."})
}

func (h *Handler) handleAppError(c *gin.Context, err error) {
	var appErr *exception.AppError

	if errors.As(err, &appErr) {
		status := http.StatusBadRequest

		switch appErr.Code {
		case exception.EntityNotFoundCode:
			status = http.StatusNotFound
		case exception.EntityConflictCode, exception.RoleAssignedCode:
			status = http.StatusConflict
		case exception.ValidationFailedCode:
			status = http.StatusUnprocessableEntity
		case exception.EmptyUpdateCode:
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(status, gin.H{
			"status":  status,
			"message": appErr.Message,
			"fields":  appErr.Fields,
		})
		return
	}

	log.Printf("[INTERNAL ERROR]: %v", err)

	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  500,
		"message": "Ocorreu um erro inesperado em nossos servidores.",
	})
}
