package models

import (
    "time"
    "database/sql"
    "encoding/json"
    // "github.com/google/uuid"
    "github.com/lib/pq"
)

type Store struct {
    ID            string     `json:"id" db:"id"`
    Title         string     `json:"title" db:"title"`
    Description   string     `json:"description" db:"description"`
    PickupTime    string     `json:"pickUpTime" db:"pickup_time"`
    Distance      *string    `json:"distance" db:"distance"`
    Price         float64    `json:"price" db:"price"`
    OriginalPrice float64    `json:"originalPrice" db:"original_price"`
    BackgroundURL string     `json:"backgroundUrl" db:"background_url"`
    AvatarURL     sql.NullString `json:"avatarUrl" db:"avatar_url"`
    ImageURL      string     `json:"imageUrl" db:"image_url"`
    Rating        float64    `json:"rating" db:"rating"`
    Reviews       int        `json:"reviews" db:"reviews"`
    Address       string     `json:"address" db:"address"`
    ItemsLeft     int        `json:"itemsLeft" db:"items_left"`
    Latitude      float64    `json:"latitude" db:"latitude"`
    Longitude     float64    `json:"longitude" db:"longitude"`
    Highlights    pq.StringArray `json:"highlights" db:"highlights"`
    IsSaved       bool       `json:"isSaved" db:"is_saved"`
    CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
    UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}

func (s Store) MarshalJSON() ([]byte, error) {
    type Alias Store // Use type alias to avoid recursion
    
    // Create a struct for JSON marshaling
    return json.Marshal(&struct {
        Alias
        AvatarURL string `json:"avatar_url,omitempty"`
    }{
        Alias: Alias(s),
        AvatarURL: s.AvatarURL.String, // This will be empty string if NULL
    })
} 