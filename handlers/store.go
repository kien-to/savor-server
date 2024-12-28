package handlers

import (
	"fmt"
	"log"
	"net/http"
	"savor-server/db"
	"savor-server/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
)

func GetStoreDetail(c *gin.Context) {
	storeID := c.Param("id")

	var modelStore models.Store
	err := db.DB.Get(&modelStore, `
		SELECT s.*, 
			   array_agg(DISTINCT sh.highlight) FILTER (WHERE sh.highlight IS NOT NULL) as highlights
		FROM stores s
		LEFT JOIN store_highlights sh ON s.id = sh.store_id
		WHERE s.id = $1
		GROUP BY s.id
	`, storeID)

	if err != nil {
		log.Println("Error fetching store:", err)
		fmt.Println("Error fetching store:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	// Check if store is saved by user
	userID := c.GetString("user_id") // From auth middleware
	var saved bool
	err = db.DB.Get(&saved, `
		SELECT EXISTS(
			SELECT 1 FROM saved_stores 
			WHERE user_id = $1 AND store_id = $2
		)
	`, userID, storeID)

	// Convert to response format
	description := ""
	if modelStore.Description.Valid {
		description = modelStore.Description.String
	}

	pickupTime := ""
	if modelStore.PickupTime.Valid {
		pickupTime = modelStore.PickupTime.String
	}

	distance := "0 km"
	if modelStore.Distance != nil {
		distance = *modelStore.Distance
	}

	phone := ""
	if modelStore.Phone.Valid {
		phone = modelStore.Phone.String
	}

	avatarURL := ""
	if modelStore.AvatarURL.Valid {
		avatarURL = modelStore.AvatarURL.String
	}

	city := ""
	if modelStore.City.Valid {
		city = modelStore.City.String
	}

	state := ""
	if modelStore.State.Valid {
		state = modelStore.State.String
	}

	zipCode := ""
	if modelStore.ZipCode.Valid {
		zipCode = modelStore.ZipCode.String
	}

	itemsLeft := 0
	if modelStore.ItemsLeft.Valid {
		itemsLeft = int(modelStore.ItemsLeft.Int64)
	}

	reviews := 0
	if modelStore.Reviews.Valid {
		reviews = int(modelStore.Reviews.Int64)
	}

	rating := 0.0
	if modelStore.Rating.Valid {
		rating = modelStore.Rating.Float64
	}

	originalPrice := 0.0
	if modelStore.OriginalPrice.Valid {
		originalPrice = modelStore.OriginalPrice.Float64
	}

	price := 0.0
	if modelStore.Price.Valid {
		price = modelStore.Price.Float64
	}

	responseStore := struct {
		ID            string         `json:"id"`
		Title         string         `json:"title"`
		Description   string         `json:"description"`
		PickupTime    string         `json:"pickUpTime"`
		Distance      string         `json:"distance"`
		Price         float64        `json:"price"`
		OriginalPrice float64        `json:"originalPrice"`
		BackgroundURL string         `json:"backgroundUrl"`
		AvatarURL     string         `json:"avatarUrl"`
		ImageURL      string         `json:"imageUrl"`
		Rating        float64        `json:"rating"`
		Reviews       int            `json:"reviews"`
		Address       string         `json:"address"`
		City          string         `json:"city"`
		State         string         `json:"state"`
		ZipCode       string         `json:"zipCode"`
		Phone         string         `json:"phone"`
		ItemsLeft     int            `json:"itemsLeft"`
		Latitude      float64        `json:"latitude"`
		Longitude     float64        `json:"longitude"`
		Highlights    pq.StringArray `json:"highlights"`
		IsSaved       bool           `json:"isSaved"`
		StoreType     string         `json:"storeType"`
		BusinessHours types.JSONText `json:"businessHours"`
	}{
		ID:            modelStore.ID,
		Title:         modelStore.Title,
		Description:   description,
		PickupTime:    pickupTime,
		Distance:      distance,
		Price:         price,
		OriginalPrice: originalPrice,
		BackgroundURL: modelStore.BackgroundURL,
		AvatarURL:     avatarURL,
		ImageURL:      modelStore.ImageURL,
		Rating:        rating,
		Reviews:       reviews,
		Address:       modelStore.Address,
		City:          city,
		State:         state,
		ZipCode:       zipCode,
		Phone:         phone,
		ItemsLeft:     itemsLeft,
		Latitude:      modelStore.Latitude,
		Longitude:     modelStore.Longitude,
		Highlights:    modelStore.Highlights,
		IsSaved:       saved,
		StoreType:     modelStore.StoreType,
		BusinessHours: modelStore.BusinessHours,
	}

	c.JSON(http.StatusOK, responseStore)
}

func ToggleSaveStore(c *gin.Context) {
	userID := c.GetString("user_id")

	// Log all headers to see if Authorization is present
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	storeID := c.Param("id")

	// First check if the store exists
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM stores WHERE id = $1)
	`, storeID)
	if err != nil || !exists {
		fmt.Println("Error fetching store:", err)
		log.Println("Error fetching store:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	// Check if store is already saved
	var isSaved bool
	err = db.DB.Get(&isSaved, `
		SELECT EXISTS(
			SELECT 1 FROM saved_stores 
			WHERE user_id = $1 AND store_id = $2
		)
	`, userID, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check save status"})
		return
	}

	// Toggle the save status
	if isSaved {
		_, err = db.DB.Exec(`
			DELETE FROM saved_stores 
			WHERE user_id = $1 AND store_id = $2
		`, userID, storeID)
	} else {
		_, err = db.DB.Exec(`
			INSERT INTO saved_stores (user_id, store_id)
			VALUES ($1, $2)
		`, userID, storeID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle save status"})
		return
	}

	// Return the new state
	var message string
	if isSaved {
		message = "Store successfully unsaved"
	} else {
		message = "Store successfully saved"
	}

	c.JSON(http.StatusOK, gin.H{
		"isSaved": !isSaved,
		"message": message,
	})
}
