import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "uk-achmadxander/models"
)

var DB *gorm.DB

func Init() {
    connStr := "postgresql://postgres:xMuGVQlTocJKWyOEizFFglIBkJJYzaoR@viaduct.proxy.rlwy.net:47282/railway"

    var err error
    DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    sqlDB, err := DB.DB()
    if err != nil {
        log.Fatal(err)
    }

    err = sqlDB.Ping()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Successfully connected to the database!")
}
