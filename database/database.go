package database

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "uk-achmadxander/models"
)

var DB *gorm.DB

func Init() {
    connStr := "postgresql://postgres:xMuGVQlTocJKWyOEizFFglIBkJJYzaoR@viaduct.proxy.rlwy.net:47282/railway"

    db, err := sql.Open("postgres", connStr)
    if err!= nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err!= nil {
        log.Fatal(err)
    }

    log.Println("Successfully connected to the database!")
}
