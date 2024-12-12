package handlers

import (
	"fmt"
	"log"
	"net/http"
	"savor-server/db"
	"savor-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveStore(c *gin.Context) {
    userID := c.GetString("user_id")
    storeID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
        return
    }

    _, err = db.DB.Exec(`
        INSERT INTO saved_stores (user_id, store_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, store_id) DO NOTHING
    `, userID, storeID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save store"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Store saved successfully"})
}

func UnsaveStore(c *gin.Context) {
    userID := c.GetString("user_id")
    storeID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
        return
    }

    result, err := db.DB.Exec(`
        DELETE FROM saved_stores
        WHERE user_id = $1 AND store_id = $2
    `, userID, storeID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsave store"})
        return
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Store not found in saved stores"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Store unsaved successfully"})
}

func GetFavorites(c *gin.Context) {
    userID := c.GetString("user_id")
    
    if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

    // fmt.Println("User ID:", userID)
    fmt.Printf("User ID: %s\n", userID)
    log.Println("User ID:", userID)

    var stores []models.Store
    err := db.DB.Select(&stores, `
        SELECT s.*, true as is_saved
        FROM stores s
        INNER JOIN saved_stores ss ON s.id = ss.store_id
        WHERE ss.user_id = $1
    `, userID)

    fmt.Println("Favorites:", stores)

    if err != nil {
        fmt.Println("Error fetching favorites:", err)
        log.Println("Error fetching favorites:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
        return
    }

    c.JSON(http.StatusOK, stores)
} 