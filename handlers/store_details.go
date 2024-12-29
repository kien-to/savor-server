package handlers

import (
	"net/http"
	"savor-server/db"
	"savor-server/models"

	"github.com/gin-gonic/gin"
)

type UpdateBagDetailsRequest struct {
	Category    string `json:"category" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Size        string `json:"size" binding:"required"`
	DailyCount  int    `json:"dailyCount" binding:"required,min=1"`
}

type UpdateScheduleRequest struct {
	Schedule []models.PickupSchedule `json:"schedule" binding:"required"`
}

func UpdateBagDetails(c *gin.Context) {
	userID := c.GetString("user_id")
	var req UpdateBagDetailsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate price based on size
	var price, minValue float64
	switch req.Size {
	case "small":
		price = 4.99
		minValue = 15.00
	case "medium":
		price = 5.99
		minValue = 18.00
	case "large":
		price = 6.99
		minValue = 21.00
	default:
		price = 5.99
		minValue = 18.00
	}

	// Update store details in a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Update store table
	_, err = tx.Exec(`
        UPDATE stores 
        SET title = $1, description = $2, price = $3, original_price = $4, 
            items_left = $5, pickup_time = $6
        WHERE owner_id = $7`,
		req.Name, req.Description, price, minValue,
		req.DailyCount, "Pick up today", userID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store"})
		return
	}

	// Also save to bag_details for reference
	_, err = tx.Exec(`
        INSERT INTO bag_details 
        (store_id, category, name, description, size, price, min_value, daily_count)
        VALUES ((SELECT id FROM stores WHERE owner_id = $1), $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (store_id) DO UPDATE
        SET category = $2, name = $3, description = $4, size = $5, 
            price = $6, min_value = $7, daily_count = $8`,
		userID, req.Category, req.Name, req.Description,
		req.Size, price, minValue, req.DailyCount)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bag details"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store and bag details updated successfully"})
}

func UpdatePickupSchedule(c *gin.Context) {
	userID := c.GetString("user_id")
	var req UpdateScheduleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get store ID
	var storeID string
	err := db.DB.QueryRow("SELECT id FROM stores WHERE owner_id = $1", userID).Scan(&storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Delete existing schedule
	_, err = tx.Exec("DELETE FROM pickup_schedules WHERE store_id = $1", storeID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
		return
	}

	// Insert new schedule
	for _, schedule := range req.Schedule {
		_, err = tx.Exec(`
            INSERT INTO pickup_schedules 
            (store_id, day, enabled, start_time, end_time)
            VALUES ($1, $2, $3, $4, $5)`,
			storeID, schedule.Day, schedule.Enabled,
			schedule.StartTime, schedule.EndTime)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully"})
}
