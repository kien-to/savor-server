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

	// Try to use environment variables first (for cloud deployment)
	// if projectID := os.Getenv("FIREBASE_PROJECT_ID"); projectID != "" {
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
			return nil, fmt.Errorf("failed to marshal Firebase config: %v", err)
		}

		// Initialize Firebase with credentials from environment
		opt := option.WithCredentialsJSON(configJSON)
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Firebase from environment: %v", err)
		}

		return app, nil
	// }

	// // Fallback to local JSON file (for local development)
	// if _, err := os.Stat("config/firebase-service-account.json"); err == nil {
	// 	opt := option.WithCredentialsFile("config/firebase-service-account.json")
	// 	app, err := firebase.NewApp(ctx, nil, opt)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to initialize Firebase from JSON file: %v", err)
	// 	}
	// 	return app, nil
	// }

	// return nil, fmt.Errorf("Firebase configuration not found. Please set environment variables or provide config/firebase-service-account.json")
}

// Helper function to get environment variable with default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
