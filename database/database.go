package database

import (
    "fmt"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "uk-achmadxander/models"
)

var DB *gorm.DB

func Init() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        panic("DATABASE_URL environment variable is not set")
    }

    fmt.Println("Connecting to database with DSN:", dsn)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect to database")
    }

    db.AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})

    DB = db
    fmt.Println("Database connected successfully")
}
