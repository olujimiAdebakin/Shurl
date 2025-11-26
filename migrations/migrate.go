package main

import (
	"log"

	"github.com/olujimiAdebakin/Shurl/initializers"
	"github.com/olujimiAdebakin/Shurl/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	log.Println("Starting database migration...")

	// Run migrations for all models
	err := initializers.DB.AutoMigrate(
		&models.User{},
		&models.Link{},
		// &models.Supplier{},
		// &models.Farmer{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("✅ Database migration completed successfully!")
}
