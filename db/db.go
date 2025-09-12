package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// dbURL = "postgresql://postgres:kientrungto9502@db.zopjihcaghguorqqopjv.supabase.co:5432/postgres"
		// dbURL = "postgresql://postgres.zopjihcaghguorqqopjv:totrungkien0905@aws-0-us-west-1.pooler.supabase.com:5432/postgres"
		dbURL = "postgres://savor_user:your_password@localhost:5432/savor?sslmode=disable"
	}

	var err error
	DB, err = sqlx.Connect("postgres", dbURL)
	// log.Printf("🔍 DEBUG: dbURL: '%s'", dbURL)
	if err != nil {
		log.Printf("error connecting to the database: %v", err)
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	return nil
}
