package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/laurencefluciano/content-api/config"
)

var globalJwks *keyfunc.JWKS

type UserData struct {
	EmailVerified bool   `json:"email_verified"`
	AuthID        string `json:"sub"`
}

func GetUserAuth(c *gin.Context) (UserData, bool) {
	val, ok := c.Get("user_auth")
	if !ok {
		return UserData{}, false
	}
	user, ok := val.(UserData)
	return user, ok
}

func AuthMiddleware() gin.HandlerFunc {
	expectedIss := config.GetEnv("SUPABASE_PROJECT_URL") + config.GetEnv("SUPABASE_AUTH_PATH")

	return func(c *gin.Context) {
		startTime := time.Now()
		tokenString := ExtractAuthorizationHeader(c)

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token ausente"})
			return
		}

		token, err := jwt.Parse(tokenString, globalJwks.Keyfunc)
		if err != nil || !token.Valid {
			log.Printf("Erro detalhado do JWT: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidas"})
			return
		}

		userMetadata, ok := claims["user_metadata"].(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidas"})
			return
		}

		emailVerified, _ := userMetadata["email_verified"].(bool)

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

		userAuth := UserData{
			EmailVerified: emailVerified,
			AuthID:        userID,
		}

		c.Set("user_auth", userAuth)

		duration := time.Since(startTime)

		log.Printf("⏱️ Tempo de execução do AuthMiddleware: %v", duration)

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

func InitJWKS() {
	startTime := time.Now()
	authURL := config.GetEnv("SUPABASE_PROJECT_URL") + config.GetEnv("SUPABASE_AUTH_PATH")
	jwksURL := authURL + "/.well-known/jwks.json"

	var err error
	// Se isso demorar, vai demorar no STARTUP do terminal, não no Postman
	globalJwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour * 24,
	})
	if err != nil {
		log.Fatalf("Erro JWKS: %v", err)
	}
	duration := time.Since(startTime)

	log.Printf("⏱️ Tempo de execução do InitJWKS: %v", duration)
}
