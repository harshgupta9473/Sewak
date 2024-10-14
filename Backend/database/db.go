package database





import (
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    connStr := os.Getenv("connStr")

    db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    
    sqlDB, err := db.DB() // Gets the underlying database/sql connection from GORM
    if err != nil {
        log.Fatal("Failed to get database connection:", err)
    }

    if err := sqlDB.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    DB = db
    log.Println("Database connection established")
}

func GetDB()*gorm.DB{
    return DB
}