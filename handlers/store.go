package handlers

import (
	"fmt"
	"log"
	"net/http"
	"savor-server/db"
	"savor-server/models"

	"github.com/gin-gonic/gin"
)

func GetStoreDetail(c *gin.Context) {
	storeID := c.Param("id")

	var store models.Store
	err := db.DB.Get(&store, `
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

	store.IsSaved = saved
	c.JSON(http.StatusOK, store)
}

func ToggleSaveStore(c *gin.Context) {
	userID := c.GetString("user_id")
    
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
		fmt.Println("Error checking save status:", err)
		log.Println("Error checking save status:", err)
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
		fmt.Println("Error toggling save status:", err)
		log.Println("Error toggling save status:", err)
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