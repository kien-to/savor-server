package config

import "github.com/stripe/stripe-go/v74"

func InitializeStripe(secretKey string) {
    stripe.Key = secretKey
} 