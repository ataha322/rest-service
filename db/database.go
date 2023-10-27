package db

import (
	"log"
	"os"
    "rest-service/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    if err := godotenv.Load(); err != nil {
        log.Fatalln("Could not load DSN")
    }
    dsn := os.Getenv("DATABASE_URL")
    log.Println("DSN loaded from .env:", dsn)

    var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Could not find the database")
	}

    log.Println("Connected to DB")
}

func AutoMigrate() {
    DB.AutoMigrate(models.Person{})
}
