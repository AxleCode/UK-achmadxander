// uk-achmadxander/database/database.go
package database

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "uk-achmadxander/models"
)

var DB *gorm.DB

func Init() {
    dsn := "host=localhost user=postgres password=watesnegoro dbname=MyGram port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect to database")
    }

    db.AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})

    DB = db
    fmt.Println("Database connected successfully")
}
