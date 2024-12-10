package db

import (
    "fmt"
    "os"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init() error {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://savor_user:your_password@localhost:5432/savor?sslmode=disable"
    }

    var err error
    DB, err = sqlx.Connect("postgres", dbURL)
    if err != nil {
        return fmt.Errorf("error connecting to the database: %v", err)
    }

    return nil
} 