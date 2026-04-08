package user

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/laurencefluciano/content-api/internal/config"
)

func IsValidAuthId(authId string) (bool, error) {
	baseURL := "https://" + config.GetEnv("SUPABASE_PROJECT_URL") + ".supabase.co/auth/v1/admin/users/"
	fullURL := baseURL + authId

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return false, fmt.Errorf("falha ao criar request: %w", err)
	}

	serviceKey := config.GetEnv("SUPABASE_API_KEY")
	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("apikey", serviceKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("falha na chamada ao Supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("erro inesperado do Supabase: status %d", resp.StatusCode)
}

func AuthMiddleware(jwksURL string) gin.HandlerFunc {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Minute * 10,
	})

	if err != nil {
		log.Fatalf("Erro ao inicializar JWKS: %v", err)
	}

	expectedIss := "https://" + config.GetEnv("SUPABASE_PROJECT_URL") + ".supabase.co/auth/v1"

	return func(c *gin.Context) {
		tokenString := ExtractAuthorizationHeader(c)

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token ausente"})
			return
		}

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidas"})
			return
		}

		if claims["iss"] != expectedIss {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Issuer não confiável"})
			return
		}

		if claims["aud"] != "authenticated" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Audience inválido"})
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Subject (sub) ausente"})
			return
		}

		c.Set("user_id", userID)
		c.Set("user_claims", claims)

		c.Next()
	}
}

func ExtractAuthorizationHeader(r *gin.Context) string {
	authHeader := r.Request.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	return strings.TrimPrefix(authHeader, "Bearer ")
}
