package handlers

import (
	"fmt"
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

type UpdateBagCountRequest struct {
	DailyCount int `json:"dailyCount" binding:"required,min=1"`
}

func UpdateBagDetails(c *gin.Context) {
	userID := c.GetString("user_id")

	// Check if this is a full update or just bag count update
	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json; type=count" {
		var req UpdateBagCountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get store ID
		var storeID string
		err := db.DB.QueryRow("SELECT id FROM stores WHERE owner_id = $1 LIMIT 1", userID).Scan(&storeID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
			return
		}

		// Update only the items_left in stores and daily_count in bag_details
		tx, err := db.DB.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
			return
		}

		// Update store table
		_, err = tx.Exec(`
			UPDATE stores 
			SET items_left = $1
			WHERE id = $2`,
			req.DailyCount, storeID)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store"})
			return
		}

		// Update bag details
		_, err = tx.Exec(`
			UPDATE bag_details 
			SET daily_count = $1
			WHERE store_id = $2`,
			req.DailyCount, storeID)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bag details"})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Bag count updated successfully"})
		return
	}

	// Handle full update
	var req UpdateBagDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First get the store ID in a separate query
	var storeID string
	err := db.DB.QueryRow("SELECT id FROM stores WHERE owner_id = $1 LIMIT 1", userID).Scan(&storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
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

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Update store table with basic details
	_, err = tx.Exec(`
        UPDATE stores 
        SET title = $1, 
            description = $2, 
            price = $3, 
            original_price = $4, 
            items_left = $5,
            store_type = $6
        WHERE id = $7`,
		req.Name,
		req.Description,
		price,
		minValue,
		req.DailyCount,
		req.Category,
		storeID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store"})
		return
	}

	// Delete existing highlights for this store
	_, err = tx.Exec(`DELETE FROM store_highlights WHERE store_id = $1`, storeID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store highlights"})
		return
	}

	// Insert new highlight for category
	_, err = tx.Exec(`
        INSERT INTO store_highlights (store_id, highlight)
        VALUES ($1, $2)`,
		storeID, req.Category)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store highlights"})
		return
	}

	// Update or insert bag details
	_, err = tx.Exec(`
        INSERT INTO bag_details 
        (store_id, category, name, description, size, price, min_value, daily_count)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (store_id) DO UPDATE
        SET category = $2, name = $3, description = $4, size = $5, 
            price = $6, min_value = $7, daily_count = $8`,
		storeID, req.Category, req.Name, req.Description,
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

	// Update store pickup time with the first enabled schedule
	enabledSchedule := ""
	for _, s := range req.Schedule {
		if s.Enabled {
			enabledSchedule = fmt.Sprintf("Pick up %s %s - %s", s.Day, s.StartTime, s.EndTime)
			break
		}
	}

	if enabledSchedule != "" {
		_, err = tx.Exec(`
            UPDATE stores 
            SET pickup_time = $1
            WHERE owner_id = $2`,
			enabledSchedule, userID)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store pickup time"})
			return
		}
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
