package envcheck

import (
	"log"

	"github.com/joho/godotenv"
)

func Init(){
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	log.Print(`All environment variables are loaded successfully!`)
}
