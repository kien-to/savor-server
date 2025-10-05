package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// FirebaseConfig represents the Firebase service account configuration
type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func InitializeFirebase() (*firebase.App, error) {
	ctx := context.Background()

	// First try to use local JSON file (for local development)
	fmt.Printf("Attempting to initialize Firebase from JSON file...\n")
	if _, err := os.Stat("config/firebase-service-account.json"); err == nil {
		fmt.Printf("Found firebase-service-account.json file\n")
		opt := option.WithCredentialsFile("config/firebase-service-account.json")
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			fmt.Printf("Failed to initialize Firebase from JSON file: %v\n", err)
			// Don't return here, try environment variables as fallback
		} else {
			fmt.Printf("Firebase initialized successfully from JSON file\n")
			return app, nil
		}
	} else {
		fmt.Printf("firebase-service-account.json not found: %v\n", err)
	}

	// Fallback to environment variables (for cloud deployment)
	fmt.Printf("Attempting to initialize Firebase with environment variables...\n")
	if projectID := os.Getenv("FIREBASE_PROJECT_ID"); projectID != "" {
		fmt.Printf("Found FIREBASE_PROJECT_ID: %s\n", projectID)

		// Check if all required environment variables are set
		requiredVars := map[string]string{
			"FIREBASE_PROJECT_ID":     os.Getenv("FIREBASE_PROJECT_ID"),
			"FIREBASE_PRIVATE_KEY_ID": os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
			"FIREBASE_PRIVATE_KEY":    os.Getenv("FIREBASE_PRIVATE_KEY"),
			"FIREBASE_CLIENT_EMAIL":   os.Getenv("FIREBASE_CLIENT_EMAIL"),
			"FIREBASE_CLIENT_ID":      os.Getenv("FIREBASE_CLIENT_ID"),
		}

		missingVars := []string{}
		for varName, value := range requiredVars {
			if value == "" {
				missingVars = append(missingVars, varName)
			}
		}

		if len(missingVars) > 0 {
			fmt.Printf("Missing required environment variables: %v\n", missingVars)
			return nil, fmt.Errorf("missing required Firebase environment variables: %v", missingVars)
		}

		// Create Firebase config from environment variables
		config := &FirebaseConfig{
			Type:                    "service_account",
			ProjectID:               os.Getenv("FIREBASE_PROJECT_ID"),
			PrivateKeyID:            os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
			PrivateKey:              os.Getenv("FIREBASE_PRIVATE_KEY"),
			ClientEmail:             os.Getenv("FIREBASE_CLIENT_EMAIL"),
			ClientID:                os.Getenv("FIREBASE_CLIENT_ID"),
			AuthURI:                 getEnvWithDefault("FIREBASE_AUTH_URI", "https://accounts.google.com/o/oauth2/auth"),
			TokenURI:                getEnvWithDefault("FIREBASE_TOKEN_URI", "https://oauth2.googleapis.com/token"),
			AuthProviderX509CertURL: getEnvWithDefault("FIREBASE_AUTH_PROVIDER_X509_CERT_URL", "https://www.googleapis.com/oauth2/v1/certs"),
			ClientX509CertURL:       os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
		}

		// Convert to JSON for Firebase SDK
		configJSON, err := json.Marshal(config)
		if err != nil {
			fmt.Printf("Failed to marshal Firebase config: %v\n", err)
			return nil, fmt.Errorf("failed to marshal Firebase config: %v", err)
		}

		// Initialize Firebase with credentials from environment
		opt := option.WithCredentialsJSON(configJSON)
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			fmt.Printf("Failed to initialize Firebase from environment: %v\n", err)
			return nil, fmt.Errorf("failed to initialize Firebase from environment: %v", err)
		}

		fmt.Printf("Firebase initialized successfully from environment variables\n")
		return app, nil
	}

	fmt.Printf("No Firebase configuration found. Please provide config/firebase-service-account.json or set environment variables\n")
	return nil, fmt.Errorf("firebase configuration not found. Please provide config/firebase-service-account.json or set environment variables")
}

// Helper function to get environment variable with default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
