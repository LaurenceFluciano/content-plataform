package main

import (
	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/user"
)

// "log"
// "github.com/google/uuid"
// "github.com/laurencefluciano/content-api/internal/config"
// "github.com/laurencefluciano/content-api/internal/database"
// "github.com/laurencefluciano/content-api/internal/user"

func main() {
	r := gin.Default()

	api := r.Group("/api/content-plataform/v1")

	userHandler := &user.Handler{}

	userHandler.InitRoutes(api)

	r.Run()
}
