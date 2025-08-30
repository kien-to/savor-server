package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SupabaseDBConfig struct {
	DirectConnection  string
	SessionPooler     string
	TransactionPooler string
	DedicatedPooler   string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
	ConnMaxIdleTime   time.Duration
}

// GetSupabaseDBConfig - Get Supabase database configuration
func GetSupabaseDBConfig() *SupabaseDBConfig {
	return &SupabaseDBConfig{
		// Direct connection (IPv6) - Best for persistent servers
		DirectConnection: os.Getenv("SUPABASE_DIRECT_URL"),

		// Session mode pooler (IPv4/IPv6) - Best for persistent servers without IPv6
		SessionPooler: os.Getenv("SUPABASE_SESSION_URL"),

		// Transaction mode pooler (IPv4/IPv6) - Best for serverless/edge functions
		TransactionPooler: os.Getenv("SUPABASE_TRANSACTION_URL"),

		// Dedicated pooler (IPv6 or with IPv4 add-on) - Best performance for paid tier
		DedicatedPooler: os.Getenv("SUPABASE_DEDICATED_URL"),

		// Connection pool settings
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
		ConnMaxIdleTime: time.Duration(getEnvAsInt("DB_CONN_MAX_IDLE_MINUTES", 5)) * time.Minute,
	}
}

// InitializeSupabaseDB - Initialize database connection using Supabase recommended methods
func InitializeSupabaseDB() (*sqlx.DB, error) {
	config := GetSupabaseDBConfig()

	// Determine which connection string to use based on availability and deployment type
	var connectionString string
	var connectionType string

	// Priority order based on Supabase recommendations:
	// 1. Dedicated pooler (best performance, paid tier)
	// 2. Direct connection (persistent servers with IPv6)
	// 3. Session pooler (persistent servers without IPv6)
	// 4. Transaction pooler (serverless/edge functions)
	// 5. Fallback to DATABASE_URL

	// Use Session Pooler (IPv4/IPv6 compatible) instead of Direct Connection (IPv6 only)
	// This fixes the "no route to host" error for environments without IPv6 support
	if config.DedicatedPooler != "" {
		connectionString = config.DedicatedPooler
		connectionType = "Dedicated Pooler"
	} else if config.SessionPooler != "" {
		connectionString = config.SessionPooler
		connectionType = "Session Pooler (IPv4/IPv6)"
	} else if config.TransactionPooler != "" {
		connectionString = config.TransactionPooler
		connectionType = "Transaction Pooler (IPv4/IPv6)"
	} else {
		// Fallback to IPv4-compatible session pooler instead of direct connection
		// connectionString = "postgresql://postgres:kientrungto9502@db.zopjihcaghguorqqopjv.supabase.co:5432/postgres"
		connectionString = os.Getenv("DATABASE_URL")
		if connectionString == "" {
			// Use IPv4-only session pooler format to avoid IPv6 connectivity issues
			// This uses the Session Pooler which supports both IPv4 and IPv6 but prefers IPv4

			connectionString = "postgresql://postgres.pdafqwrgdgbqbgbbtqqa:totrungkien0905@aws-1-us-east-1.pooler.supabase.com:5432/postgres"
			// "postgresql://postgres.zopjihcaghguorqqopjv:totrungkien0905@aws-0-us-west-1.pooler.supabase.com:5432/postgres"
		}
		connectionType = "Fallback Session Pooler (IPv4)"
	}

	log.Printf("🔗 Connecting to Supabase using: %s", connectionType)
	log.Printf("📡 Connection string: %s", connectionString)

	// Create database connection
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Printf("🔍 DEBUG: error connecting to Supabase database: %v", err)
		return nil, fmt.Errorf("failed to connect to Supabase database: %v", err)
	}

	// Configure connection pool for optimal performance
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Printf("🔍 DEBUG: error pinging Supabase database: %v", err)
		return nil, fmt.Errorf("failed to ping Supabase database: %v", err)
	}

	log.Printf("✅ Successfully connected to Supabase PostgreSQL")
	// log.Printf("📊 Connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v",
	// 	config.MaxOpenConns, config.MaxIdleConns, config.ConnMaxLifetime, config.ConnMaxIdleTime)

	return db, nil
}

// Helper function to get environment variable as integer with default
func getEnvAsInt(name string, defaultVal int) int {
	if str := os.Getenv(name); str != "" {
		if val, err := strconv.Atoi(str); err == nil {
			return val
		}
	}
	return defaultVal
}
