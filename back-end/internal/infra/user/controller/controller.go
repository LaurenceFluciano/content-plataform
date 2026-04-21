package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/app/user/usecase"
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/infra/middleware"
)

type UserController struct {
	registerUser   usecase.RegisterUserUseCase
	updateProfile  usecase.UpdateProfileUseCase
	becameProducer usecase.BecomeProducerUserCase
	repo           user.Repository
}

func CreateController(r *gin.RouterGroup, repo user.Repository) *UserController {
	handler := &UserController{
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

func (h *UserController) InitRoutes(r *gin.RouterGroup) {
	auth := middleware.AuthMiddleware()

	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.LatencyLogger())
	userRoutes.Use(auth)
	{
		userRoutes.POST("/", h.register)
		userRoutes.GET("/me/profile", h.getMe)
		userRoutes.PATCH("/me/profile", h.updateMe)
		userRoutes.POST("/become-producer", h.becomeProducer)
	}
}
