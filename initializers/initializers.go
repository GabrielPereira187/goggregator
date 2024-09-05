package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadInitializers() {
	err := godotenv.Load("/home/gabriel/Documentos/boot.dev/golang/blog-aggregator-api/.env")

	if err != nil {
		log.Fatal("erro")
	}
}
