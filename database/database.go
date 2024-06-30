import (
    "fmt"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "uk-achmadxander/models"
)

var DB *gorm.DB

func Init() {
    host := os.Getenv("PGHOST")
    user := os.Getenv("PGUSER")
    password := os.Getenv("PGPASSWORD")
    dbname := os.Getenv("PGDATABASE")
    port := os.Getenv("PGPORT")

    fmt.Println("Environment variables:")
    fmt.Println("PGHOST:", host)
    fmt.Println("PGUSER:", user)
    fmt.Println("PGPASSWORD:", password)
    fmt.Println("PGDATABASE:", dbname)
    fmt.Println("PGPORT:", port)

    if host == "" || user == "" || password == "" || dbname == "" || port == "" {
        panic("One or more environment variables are not set")
    }

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        host, user, password, dbname, port,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect to database")
    }

    db.AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})

    DB = db
    fmt.Println("Database connected successfully")
}
