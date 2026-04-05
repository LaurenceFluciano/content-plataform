package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/laurencefluciano/content-api/internal/config"
	"github.com/laurencefluciano/content-api/internal/database"
	"github.com/laurencefluciano/content-api/internal/user"
)

func main() {
	config.InitEnv()
	db, err := database.NewPostgresDB()

	if err != nil {
		log.Fatal("Error ao efetuar a conexão com o banco de dados: ", err)
	}

	repo := user.NewRepository(db)

	repo.CreateUser("Example", uuid.NewString(), user.Active)

}
