package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"savor-server/db"
	"time"

	"github.com/gin-gonic/gin"
)

type StoreOwnerReservation struct {
	ID              string     `json:"id"`
	CustomerName    string     `json:"customerName"`
	CustomerEmail   string     `json:"customerEmail"`
	PhoneNumber     string     `json:"phoneNumber"`
	Quantity        int        `json:"quantity"`
	TotalAmount     float64    `json:"totalAmount"`
	Status          string     `json:"status"`
	PickupTime      *string    `json:"pickupTime,omitempty"`
	PickupTimestamp *time.Time `json:"pickupTimestamp,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	StoreName       string     `json:"storeName"`
	StoreImage      string     `json:"storeImage"`
	StoreAddress    string     `json:"storeAddress"`
}

type StoreOwnerSettings struct {
	SurpriseBoxes int     `json:"surpriseBoxes"`
	Price         float64 `json:"price"`
	IsSelling     bool    `json:"isSelling"`
}

type UpdateReservationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type UpdateStoreSettingsRequest struct {
	SurpriseBoxes int     `json:"surpriseBoxes"`
	Price         float64 `json:"price"`
	IsSelling     bool    `json:"isSelling"`
}

// GetStoreOwnerReservations gets all reservations for a store owner's store
func GetStoreOwnerReservations(c *gin.Context) {
	fmt.Printf("DEBUG: GetStoreOwnerReservations called\n")
	userID := c.GetString("user_id")
	fmt.Printf("DEBUG: Retrieved userID from context: '%s'\n", userID)

	if userID == "" {
		fmt.Printf("ERROR: User not authenticated - userID is empty\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// First, get the store ID for this user
	fmt.Printf("DEBUG: Looking for store with owner_id: %s\n", userID)
	var storeID string
	err := db.DB.QueryRow(`
		SELECT id FROM stores WHERE owner_id = $1
	`, userID).Scan(&storeID)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("WARN: No store found for userID: %s - returning empty reservations\n", userID)
			c.JSON(http.StatusOK, gin.H{
				"currentReservations": []StoreOwnerReservation{},
				"pastReservations":    []StoreOwnerReservation{},
				"currentCount":        0,
				"pastCount":           0,
			})
			return
		}
		fmt.Printf("ERROR: Failed to query store for userID %s: %v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get store"})
		return
	}

	fmt.Printf("DEBUG: Found store with ID: %s for userID: %s\n", storeID, userID)

	// Get all reservations for this store
	fmt.Printf("DEBUG: Querying reservations for store_id: %s\n", storeID)
	rows, err := db.DB.Query(`
		SELECT 
			r.id,
			COALESCE(r.customer_name, u.email, 'Guest User') as customer_name,
			COALESCE(r.customer_email, u.email, '') as customer_email,
			COALESCE(r.phone_number, '') as phone_number,
			r.quantity,
			r.total_amount,
			r.status,
			r.pickup_time,
			r.pickup_timestamp,
			r.created_at,
			s.title as store_name,
			s.image_url as store_image,
			s.address as store_address
		FROM reservations r
		LEFT JOIN users u ON r.user_id = u.id::text
		JOIN stores s ON r.store_id = s.id
		WHERE r.store_id = $1 
		ORDER BY r.pickup_timestamp DESC, r.created_at DESC
	`, storeID)

	if err != nil {
		fmt.Printf("ERROR: Failed to query reservations for store_id %s: %v\n", storeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reservations"})
		return
	}
	defer rows.Close()

	var currentReservations []StoreOwnerReservation
	var pastReservations []StoreOwnerReservation
	now := time.Now()
	twentyFourHoursAgo := now.Add(-24 * time.Hour)

	for rows.Next() {
		var res StoreOwnerReservation
		err := rows.Scan(
			&res.ID,
			&res.CustomerName,
			&res.CustomerEmail,
			&res.PhoneNumber,
			&res.Quantity,
			&res.TotalAmount,
			&res.Status,
			&res.PickupTime,
			&res.PickupTimestamp,
			&res.CreatedAt,
			&res.StoreName,
			&res.StoreImage,
			&res.StoreAddress,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan reservation"})
			return
		}

		// Categorize based on 24-hour window
		if res.CreatedAt.After(twentyFourHoursAgo) {
			// Created within last 24 hours - current reservation
			currentReservations = append(currentReservations, res)
		} else {
			// Created more than 24 hours ago - past reservation
			pastReservations = append(pastReservations, res)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"currentReservations": currentReservations,
		"pastReservations":    pastReservations,
		"currentCount":        len(currentReservations),
		"pastCount":           len(pastReservations),
	})
}

// UpdateReservationStatus updates the status of a reservation
func UpdateReservationStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	reservationID := c.Param("id")
	if reservationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reservation ID is required"})
		return
	}

	var req UpdateReservationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	if req.Status != "active" && req.Status != "picked_up" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be 'active' or 'picked_up'"})
		return
	}

	// First, verify that the reservation belongs to the user's store
	var storeID string
	err := db.DB.QueryRow(`
		SELECT s.id 
		FROM stores s
		JOIN reservations r ON s.id = r.store_id
		WHERE s.owner_id = $1 AND r.id = $2
	`, userID, reservationID).Scan(&storeID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found or not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify reservation"})
		return
	}

	// Update the reservation status
	_, err = db.DB.Exec(`
		UPDATE reservations 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`, req.Status, reservationID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reservation status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation status updated successfully",
		"status":  req.Status,
	})
}

// GetStoreOwnerSettings gets the current store settings
func GetStoreOwnerSettings(c *gin.Context) {
	fmt.Printf("DEBUG: GetStoreOwnerSettings called")
	userID := c.GetString("user_id")
	fmt.Printf("DEBUG: Retrieved userID from context: '%s'", userID)

	if userID == "" {
		fmt.Printf("ERROR: User not authenticated - userID is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var settings StoreOwnerSettings
	fmt.Printf("DEBUG: Querying store settings for userID: %s", userID)
	err := db.DB.QueryRow(`
		SELECT 
			COALESCE(bags_available, 10) as surprise_boxes,
			COALESCE(price, 5.0) as price,
			COALESCE(is_selling, false) as is_selling
		FROM stores 
		WHERE owner_id = $1
	`, userID).Scan(&settings.SurpriseBoxes, &settings.Price, &settings.IsSelling)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("WARN: No store found for userID: %s - returning default settings\n", userID)
			settings = StoreOwnerSettings{SurpriseBoxes: 10, Price: 5.0, IsSelling: false}
		} else {
			fmt.Printf("ERROR: Failed to query store settings for userID %s: %v\n", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get store settings"})
			return
		}
	}

	fmt.Printf("DEBUG: Responding store settings for userID %s: %+v\n", userID, settings)

	c.JSON(http.StatusOK, settings)
}

// UpdateStoreOwnerSettings updates the store settings
func UpdateStoreOwnerSettings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req UpdateStoreSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if req.SurpriseBoxes < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Surprise boxes count cannot be negative"})
		return
	}
	if req.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price cannot be negative"})
		return
	}

	// Update store settings
	_, err := db.DB.Exec(`
		UPDATE stores 
		SET 
			bags_available = $1,
			price = $2,
			is_selling = $3,
			updated_at = NOW()
		WHERE owner_id = $4
	`, req.SurpriseBoxes, req.Price, req.IsSelling, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Store settings updated successfully",
		"settings": StoreOwnerSettings{
			SurpriseBoxes: req.SurpriseBoxes,
			Price:         req.Price,
			IsSelling:     req.IsSelling,
		},
	})
}

// GetStoreOwnerStats gets statistics for the store owner
func GetStoreOwnerStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get store ID
	var storeID string
	err := db.DB.QueryRow(`
		SELECT id FROM stores WHERE owner_id = $1
	`, userID).Scan(&storeID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get store"})
		return
	}

	// Get stats for current and past reservations
	now := time.Now()

	var currentStats struct {
		TotalReservations    int     `json:"totalReservations"`
		ActiveReservations   int     `json:"activeReservations"`
		PickedUpReservations int     `json:"pickedUpReservations"`
		TotalRevenue         float64 `json:"totalRevenue"`
	}

	var pastStats struct {
		TotalReservations    int     `json:"totalReservations"`
		ActiveReservations   int     `json:"activeReservations"`
		PickedUpReservations int     `json:"pickedUpReservations"`
		TotalRevenue         float64 `json:"totalRevenue"`
	}

	// Get current reservations stats (future pickup times)
	err = db.DB.QueryRow(`
		SELECT 
			COUNT(*) as total_reservations,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_reservations,
			COUNT(CASE WHEN status = 'picked_up' THEN 1 END) as picked_up_reservations,
			COALESCE(SUM(total_amount), 0) as total_revenue
		FROM reservations 
		WHERE store_id = $1 
		AND pickup_timestamp IS NOT NULL 
		AND pickup_timestamp > $2
	`, storeID, now).Scan(
		&currentStats.TotalReservations,
		&currentStats.ActiveReservations,
		&currentStats.PickedUpReservations,
		&currentStats.TotalRevenue,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current stats"})
		return
	}

	// Get past reservations stats (past pickup times or no pickup_timestamp)
	err = db.DB.QueryRow(`
		SELECT 
			COUNT(*) as total_reservations,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_reservations,
			COUNT(CASE WHEN status = 'picked_up' THEN 1 END) as picked_up_reservations,
			COALESCE(SUM(total_amount), 0) as total_revenue
		FROM reservations 
		WHERE store_id = $1 
		AND (pickup_timestamp IS NULL OR pickup_timestamp <= $2)
	`, storeID, now).Scan(
		&pastStats.TotalReservations,
		&pastStats.ActiveReservations,
		&pastStats.PickedUpReservations,
		&pastStats.TotalRevenue,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get past stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"current": currentStats,
		"past":    pastStats,
		"date":    now.Format("2006-01-02"),
	})
}
