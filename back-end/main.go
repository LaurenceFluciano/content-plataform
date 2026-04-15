package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/config"
	"github.com/laurencefluciano/content-api/internal/infra/database"
	"github.com/laurencefluciano/content-api/internal/infra/http"
)

func main() {
	/* --- Env Config --- */
	config.InitEnv()

	/* --- Database Init Config --- */
	db, err := database.NewPostgresDB()

	if err != nil {
		log.Fatal("[ ERR ] Ocorreu um erro ao tentar conectar com o banco de dados: ", err)
		return
	}

	if config.GetEnv("RUN_MIGRATIONS") == "true" {
		err = db.AutoMigrate(
			database.GetModels()...,
		)
	}

	err = database.CreateRoles(db)

	if err != nil {
		log.Fatal("[ ERR ] Ocorreu um erro ao tentar conectar com o banco de dados: ", err)
		return
	}

	userRepository := database.NewGormUserRepository(db)

	/* --- Routes Init Config --- */
	routes := gin.Default()

	api := routes.Group("/api/content-plataform/v1")

	http.CreateHandler(api, userRepository)

	routes.Run(":8080")
}
