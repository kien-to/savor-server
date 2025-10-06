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
	PickupTimestamp sql.NullTime    `json:"pickupTimestamp" db:"pickup_timestamp"`
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
	StoreType       sql.NullString  `json:"storeType" db:"store_type"`
	BusinessHours   types.JSONText  `json:"businessHours" db:"business_hours"`
	CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time       `json:"updatedAt" db:"updated_at"`
}

func (s Store) MarshalJSON() ([]byte, error) {
	// Create a struct for JSON marshaling with proper null handling
	result := struct {
		ID              string     `json:"id"`
		OwnerID         string     `json:"ownerId"`
		Title           string     `json:"title"`
		Description     *string    `json:"description"`
		PickupTime      *string    `json:"pickUpTime"`
		PickupTimestamp *time.Time `json:"pickupTimestamp"`
		Distance        *string    `json:"distance"`
		Price           *float64   `json:"price"`
		OriginalPrice   *float64   `json:"originalPrice"`
		DiscountedPrice *float64   `json:"discountedPrice"`
		BackgroundURL   string     `json:"backgroundUrl"`
		AvatarURL       *string    `json:"avatarUrl"`
		ImageURL        string     `json:"imageUrl"`
		Rating          *float64   `json:"rating"`
		Reviews         *int64     `json:"reviews"`
		ReviewsCount    *int64     `json:"reviewsCount"`
		Address         string     `json:"address"`
		City            *string    `json:"city"`
		State           *string    `json:"state"`
		ZipCode         *string    `json:"zipCode"`
		Country         *string    `json:"country"`
		Phone           *string    `json:"phone"`
		ItemsLeft       *int64     `json:"itemsLeft"`
		BagsAvailable   *int64     `json:"bagsAvailable"`
		Latitude        float64    `json:"latitude"`
		Longitude       float64    `json:"longitude"`
		GoogleMapsURL   *string    `json:"googleMapsUrl"`
		Highlights      []string   `json:"highlights"`
		IsSaved         bool       `json:"isSaved"`
		IsSelling       bool       `json:"isSelling"`
		StoreType       *string    `json:"storeType"`
		CreatedAt       time.Time  `json:"createdAt"`
		UpdatedAt       time.Time  `json:"updatedAt"`
	}{
		ID:            s.ID,
		OwnerID:       s.OwnerID,
		Title:         s.Title,
		BackgroundURL: s.BackgroundURL,
		ImageURL:      s.ImageURL,
		Address:       s.Address,
		Latitude:      s.Latitude,
		Longitude:     s.Longitude,
		Highlights:    s.Highlights,
		IsSaved:       s.IsSaved,
		IsSelling:     s.IsSelling,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}

	// Handle nullable strings
	if s.Description.Valid {
		result.Description = &s.Description.String
	}
	if s.PickupTime.Valid {
		result.PickupTime = &s.PickupTime.String
	}
	if s.PickupTimestamp.Valid {
		result.PickupTimestamp = &s.PickupTimestamp.Time
	}
	if s.AvatarURL.Valid {
		result.AvatarURL = &s.AvatarURL.String
	}
	if s.City.Valid {
		result.City = &s.City.String
	}
	if s.State.Valid {
		result.State = &s.State.String
	}
	if s.ZipCode.Valid {
		result.ZipCode = &s.ZipCode.String
	}
	if s.Country.Valid {
		result.Country = &s.Country.String
	}
	if s.Phone.Valid {
		result.Phone = &s.Phone.String
	}
	if s.GoogleMapsURL.Valid {
		result.GoogleMapsURL = &s.GoogleMapsURL.String
	}
	if s.StoreType.Valid {
		result.StoreType = &s.StoreType.String
	}

	// Handle nullable floats
	if s.Price.Valid {
		result.Price = &s.Price.Float64
	}
	if s.OriginalPrice.Valid {
		result.OriginalPrice = &s.OriginalPrice.Float64
	}
	if s.DiscountedPrice.Valid {
		result.DiscountedPrice = &s.DiscountedPrice.Float64
	}
	if s.Rating.Valid {
		result.Rating = &s.Rating.Float64
	}

	// Handle nullable ints
	if s.Reviews.Valid {
		result.Reviews = &s.Reviews.Int64
	}
	if s.ReviewsCount.Valid {
		result.ReviewsCount = &s.ReviewsCount.Int64
	}
	if s.ItemsLeft.Valid {
		result.ItemsLeft = &s.ItemsLeft.Int64
	}
	if s.BagsAvailable.Valid {
		result.BagsAvailable = &s.BagsAvailable.Int64
	}

	// Handle distance
	result.Distance = s.Distance

	return json.Marshal(result)
}

func (s *Store) GetBusinessHours() ([]BusinessHours, error) {
	if len(s.BusinessHours) == 0 {
		return []BusinessHours{}, nil
	}
	var hours []BusinessHours
	err := s.BusinessHours.Unmarshal(&hours)
	return hours, err
}
