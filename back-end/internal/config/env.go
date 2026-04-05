package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Config: Arquivo .env não encontrado, usando variáveis de sistema.")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
