package main

import (
	"go-data-migration/models"
	"go-data-migration/services"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting data migration process...")

	// Read configuration
	configData, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	var config models.Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal("Error parsing config:", err)
	}

	log.Printf("Configuration loaded successfully. Source DB: %s, Destination DB: %s",
		config.Source.Database,
		config.Destination.Database)

	// Connect to databases
	sourceDB, err := services.ConnectDB(config.Source)
	if err != nil {
		log.Fatal("Error connecting to source database:", err)
	}
	defer sourceDB.Close()
	log.Println("Connected to source database successfully")

	destDB, err := services.ConnectDB(config.Destination)
	if err != nil {
		log.Fatal("Error connecting to destination database:", err)
	}
	defer destDB.Close()
	log.Println("Connected to destination database successfully")

	// Create migration service
	migrationService := services.NewMigrationService(sourceDB, destDB, config.Migration)
	log.Println("Starting migration process...")

	// Perform migration
	err = migrationService.MigrateData()
	if err != nil {
		log.Fatal("Error during migration:", err)
	}

	log.Println("Migration completed successfully!")
}
