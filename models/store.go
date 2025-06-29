package models

import (
	"database/sql"
	"encoding/json"
	"time"

	// "github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
)

type BusinessHours struct {
	Day       string `json:"day"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Enabled   bool   `json:"enabled"`
}

type Store struct {
	ID              string          `json:"id" db:"id"`
	OwnerID         string          `json:"ownerId" db:"owner_id"`
	Title           string          `json:"title" db:"title"`
	Description     sql.NullString  `json:"description" db:"description"`
	PickupTime      sql.NullString  `json:"pickUpTime" db:"pickup_time"`
	Distance        *string         `json:"distance" db:"distance"`
	Price           sql.NullFloat64 `json:"price" db:"price"`
	OriginalPrice   sql.NullFloat64 `json:"originalPrice" db:"original_price"`
	DiscountedPrice sql.NullFloat64 `json:"discountedPrice" db:"discounted_price"`
	BackgroundURL   string          `json:"backgroundUrl" db:"background_url"`
	AvatarURL       sql.NullString  `json:"avatarUrl" db:"avatar_url"`
	ImageURL        string          `json:"imageUrl" db:"image_url"`
	Rating          sql.NullFloat64 `json:"rating" db:"rating"`
	Reviews         sql.NullInt64   `json:"reviews" db:"reviews"`
	ReviewsCount    sql.NullInt64   `json:"reviewsCount" db:"reviews_count"`
	Address         string          `json:"address" db:"address"`
	City            sql.NullString  `json:"city" db:"city"`
	State           sql.NullString  `json:"state" db:"state"`
	ZipCode         sql.NullString  `json:"zipCode" db:"zip_code"`
	Country         sql.NullString  `json:"country" db:"country"`
	Phone           sql.NullString  `json:"phone" db:"phone"`
	ItemsLeft       sql.NullInt64   `json:"itemsLeft" db:"items_left"`
	BagsAvailable   sql.NullInt64   `json:"bagsAvailable" db:"bags_available"`
	Latitude        float64         `json:"latitude" db:"latitude"`
	Longitude       float64         `json:"longitude" db:"longitude"`
	GoogleMapsURL   sql.NullString  `json:"googleMapsUrl" db:"google_maps_url"`
	Highlights      pq.StringArray  `json:"highlights" db:"highlights"`
	IsSaved         bool            `json:"isSaved" db:"is_saved"`
	IsSelling       bool            `json:"isSelling" db:"is_selling"`
	StoreType       string          `json:"storeType" db:"store_type"`
	BusinessHours   types.JSONText  `json:"businessHours" db:"business_hours"`
	CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time       `json:"updatedAt" db:"updated_at"`
}

func (s Store) MarshalJSON() ([]byte, error) {
	type Alias Store // Use type alias to avoid recursion

	// Create a struct for JSON marshaling
	return json.Marshal(&struct {
		Alias
		AvatarURL string `json:"avatar_url,omitempty"`
	}{
		Alias:     Alias(s),
		AvatarURL: s.AvatarURL.String, // This will be empty string if NULL
	})
}

func (s *Store) GetBusinessHours() ([]BusinessHours, error) {
	if len(s.BusinessHours) == 0 {
		return []BusinessHours{}, nil
	}
	var hours []BusinessHours
	err := s.BusinessHours.Unmarshal(&hours)
	return hours, err
}
