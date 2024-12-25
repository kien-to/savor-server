package handlers

import (
	"net/http"
	"savor-server/db"
	"savor-server/models"

	"github.com/gin-gonic/gin"
)

type BusinessDetails struct {
	StoreName string  `json:"businessName" binding:"required"`
	StoreType string  `json:"storeType" binding:"required"`
	Street    string  `json:"street" binding:"required"`
	City      string  `json:"city" binding:"required"`
	State     string  `json:"state" binding:"required"`
	ZipCode   string  `json:"zipCode" binding:"required"`
	Country   string  `json:"country" binding:"required"`
	Phone     string  `json:"phone" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func CreateStore(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var details BusinessDetails
	if err := c.ShouldBindJSON(&details); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	var storeID string
	err = tx.QueryRow(`
        INSERT INTO stores (
            owner_id, title, store_type, address, city, 
            state, zip_code, country, phone, latitude, longitude
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id`,
		userID, details.StoreName, details.StoreType, details.Street,
		details.City, details.State, details.ZipCode, details.Country,
		details.Phone, details.Latitude, details.Longitude,
	).Scan(&storeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Store created successfully",
		"id":      storeID,
	})
}

func GetMyStore(c *gin.Context) {
	userID := c.GetString("user_id")

	var store models.Store
	err := db.DB.Get(&store, `
        SELECT * FROM stores WHERE owner_id = $1
    `, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	c.JSON(http.StatusOK, store)
}

func UpdateStore(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var store models.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
        UPDATE stores 
        SET title = $1, description = $2, address = $3, city = $4, 
            state = $5, zip_code = $6, phone = $7, store_type = $8,
            latitude = $9, longitude = $10
        WHERE owner_id = $11
        RETURNING id`

	var storeID string
	err := db.DB.QueryRow(
		query,
		store.Title, store.Description, store.Address, store.City,
		store.State, store.ZipCode, store.Phone, store.StoreType,
		store.Latitude, store.Longitude, userID,
	).Scan(&storeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": storeID})
}
