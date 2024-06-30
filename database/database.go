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
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("PGHOST"),
        os.Getenv("PGUSER"),
        os.Getenv("PGPASSWORD"),
        os.Getenv("PGDATABASE"),
        os.Getenv("PGPORT"),
    )
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect to database")
    }

    db.AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})

    DB = db
    fmt.Println("Database connected successfully")
}
