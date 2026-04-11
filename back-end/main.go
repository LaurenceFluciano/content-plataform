package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/config"
	"github.com/laurencefluciano/content-api/internal/database"
	"github.com/laurencefluciano/content-api/internal/user"
)

// "log"
// "github.com/google/uuid"
// "github.com/laurencefluciano/content-api/internal/user"

func main() {
	/* --- Env Config --- */
	config.InitEnv()

	/* --- Database Init Config --- */
	db, err := database.NewPostgresDB()

	if err != nil {
		log.Fatal("[ ERR ] Ocorreu um erro ao tentar conectar com o banco de dados: ", err)
		return
	}

	err = db.AutoMigrate(
		user.GetModels()...,
	)

	if err != nil {
		log.Fatal("[ ERR ] Ocorreu um erro ao tentar conectar com o banco de dados: ", err)
		return
	}

	user.NewRepository(db)

	/* --- Routes Init Config --- */
	routes := gin.Default()

	api := routes.Group("/api/content-plataform/v1")

	user.CreateHandler(api)

	routes.Run(":8080")
}
