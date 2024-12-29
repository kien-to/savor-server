package models

import "time"

type BagDetails struct {
	ID          string    `json:"id" db:"id"`
	StoreID     string    `json:"storeId" db:"store_id"`
	Category    string    `json:"category" db:"category"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Size        string    `json:"size" db:"size"`
	Price       float64   `json:"price" db:"price"`
	MinValue    float64   `json:"minValue" db:"min_value"`
	DailyCount  int       `json:"dailyCount" db:"daily_count"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type PickupSchedule struct {
	ID        string    `json:"id" db:"id"`
	StoreID   string    `json:"storeId" db:"store_id"`
	Day       string    `json:"day" db:"day"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	StartTime string    `json:"startTime" db:"start_time"`
	EndTime   string    `json:"endTime" db:"end_time"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
