package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	// Load .env file for local development
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, using environment variables from system")
	}

	// Verify required environment variables
	fmt.Printf("POSTGRES_HOST: '%s'\n", os.Getenv("POSTGRES_HOST"))
	fmt.Printf("POSTGRES_USER: '%s'\n", os.Getenv("POSTGRES_USER"))
	fmt.Printf("POSTGRES_DB: '%s'\n", os.Getenv("POSTGRES_DB"))
}
