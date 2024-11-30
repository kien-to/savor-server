package config

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func InitializeFirebase() (*firebase.App, error) {
	ctx := context.Background()
	
	// Make sure to place your Firebase service account key JSON file in the config directory
	opt := option.WithCredentialsFile("config/firebase-service-account.json")
	
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	
	return app, nil
} 