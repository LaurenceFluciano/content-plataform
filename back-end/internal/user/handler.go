package user

/*

*

* ADICINAR LIMITE DE 10 URLS de WEBSITES DO PERFIL DE PRODUTOR

*

 */

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func isValidStoragePath(path string) bool {
	if path == "" {
		return false
	}

	if !strings.HasPrefix(path, "avatars/") {
		return false
	}

	ext := strings.ToLower(filepath.Ext(path))
	validExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}

	return validExtensions[ext]
}

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
		userRoutes.PUT("/me", AuthMiddleware(), h.upsert)
		userRoutes.GET("/me/profile", AuthMiddleware(), h.getMe)
		userRoutes.PATCH("/me/profile", AuthMiddleware(), h.updateMe)
		userRoutes.POST("/become-producer", AuthMiddleware(), h.becomeProducer)
	}
}

// Nome não esta sendo atualizado
// Mostrar caso o campo atualize
func (h *Handler) upsert(c *gin.Context) {
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
		name = strings.ToLower(*dto.Name)
	}

	userClaims, ok := GetUserAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	if !userClaims.EmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email não confirmado"})
		return
	}

	user, err := h.repo.GetUserByAuthId(userClaims.AuthID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err := h.repo.RegisterNewUser(&SaveUserParams{
				name:   name,
				authId: userClaims.AuthID,
				status: Active,
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"id": user.ID, "status": "created"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro de conexão com o banco"})
		return
	}

	if len(name) > 1 {
		err, existingUser := h.repo.GetUserByName(name)

		if err == nil {
			if existingUser.AuthId != userClaims.AuthID {
				c.JSON(http.StatusConflict, gin.H{"error": "Nome de usuário já está sendo utilizado."})
				return
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao validar disponibilidade do nome"})
			return
		} else {
			err = h.repo.UpdateProfileByAuthId(&UpdateUserProfile{
				name: name,
			}, userClaims.AuthID)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar perfil"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     user.ID,
		"status": "exists",
	})
}

func (h *Handler) getMe(c *gin.Context) {
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
	var dto UpdateUserRequest

	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

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

	allowedFields := GetEditableFieldsMap(profile.User.Roles)
	userUpdate := &UpdateUserProfile{}

	var fields_errors []string

	if allowedFields["Name"] {

		if dto.Name != nil {
			if len(*dto.Name) < 2 {
				fields_errors = append(
					fields_errors,
					"Nome deve ter no minimo 2 caracteres.",
				)
			} else {
				userUpdate.name = strings.ToLower(*dto.Name)

				err, existingUser := h.repo.GetUserByName(userUpdate.name)

				if err == nil {
					if existingUser.AuthId != profile.User.AuthId {
						fields_errors = append(fields_errors, "Este nome de usuário já está sendo utilizado.")
					}
				} else if !errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
					return
				}

			}
		}

	}

	if allowedFields["AvatarUrl"] {

		if dto.AvatarUrl != nil {
			if !isValidStoragePath(*dto.AvatarUrl) {
				fields_errors = append(
					fields_errors,
					"URL da imagem do Avatar é invalida, verifique se o tipo da imagem é jpg|png|jpeg.",
				)
			} else {
				userUpdate.avatarUrl = *dto.AvatarUrl
			}

		}

	}

	if allowedFields["Bio"] {

		if dto.Bio != nil {
			if len(*dto.Bio) < 10 {
				fields_errors = append(
					fields_errors,
					"A bio fornecida é invalida, ela deve ter pelo menos 10 caracteres.",
				)
			} else {
				userUpdate.bio = *dto.Bio
			}
		}

	}

	if allowedFields["Websites"] {
		if dto.Websites != nil {
			hasErrors := false
			for _, url := range *dto.Websites {
				if len(url) < 5 {
					fields_errors = append(fields_errors, "URL inválida no portfólio")
					hasErrors = true
					break
				}
			}
			if !hasErrors {
				userUpdate.websites = append(profile.ProducerProfile.Websites, *dto.Websites...)
			}
		}
	}

	if allowedFields["Social"] {
		if dto.Social != nil {
			cleanSocial := SocialLinks{}

			if dto.Social.Facebook != "" {
				cleanSocial.Facebook = strings.TrimSpace(dto.Social.Facebook)
			}
			if dto.Social.Youtube != "" {
				cleanSocial.Youtube = strings.TrimSpace(dto.Social.Youtube)
			}
			if dto.Social.Instagram != "" {
				cleanSocial.Instagram = strings.TrimSpace(dto.Social.Instagram)
			}
			if dto.Social.Tiktok != "" {
				cleanSocial.Tiktok = strings.TrimSpace(dto.Social.Tiktok)
			}
			if dto.Social.X != "" {
				cleanSocial.X = strings.TrimSpace(dto.Social.X)
			}

			userUpdate.social = cleanSocial
		}
	}

	if len(fields_errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Campos inválidos",
			"fields": fields_errors,
		})
		return
	}

	if userUpdate.name == "" &&
		userUpdate.avatarUrl == "" &&
		userUpdate.bio == "" &&
		len(userUpdate.websites) == 0 &&
		userUpdate.social == (SocialLinks{}) {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Nada para atualizar"})
		return
	}

	err = h.repo.UpdateProfileByAuthId(userUpdate, profile.User.AuthId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar perfil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil de usuário atualizado."})
}

func (h *Handler) becomeProducer(c *gin.Context) {
	userClaims, ok := GetUserAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	user, err := h.repo.GetUserByAuthId(userClaims.AuthID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error."})
		return
	}

	if user.HasRole(RoleProducer) {
		c.JSON(http.StatusConflict, gin.H{"message": "Você já é um produtor."})
		return
	}

	err = h.repo.InitializeProducerProfile(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar upgrade de conta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parabéns! Você agora é um produtor."})
}
