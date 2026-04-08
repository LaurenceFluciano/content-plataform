package user

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/laurencefluciano/content-api/internal/config"
)

func AuthMiddleware(jwksURL string) func(http.Handler) http.Handler {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Minute * 10,
	})

	if err != nil {
		log.Fatalf("Erro ao pegar JWKS: %v", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token ausente", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, jwks.Keyfunc)
			if err != nil || !token.Valid {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Claims inválidas", http.StatusUnauthorized)
				return
			}

			expectedIss := "https://" + config.GetEnv("SUPABASE_PROJECT_URL") + ".supabase.co/auth/v1"
			if claims["iss"] != expectedIss {
				http.Error(w, "Issuer inválido", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Usuário inválido", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
