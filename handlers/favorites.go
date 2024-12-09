package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// In-memory storage for saved stores (in production, this should be in a database)
var savedStores = make(map[string][]string) // userID -> []storeID

var mockHomeResponse = struct {
    UserLocation struct {
        City     string  `json:"city"`
        Distance float64 `json:"distance"`
    } `json:"userLocation"`
    PickUpToday    []Store `json:"pickUpToday"`
    PickUpTomorrow []Store `json:"pickUpTomorrow"`
    EmailVerified  bool    `json:"emailVerified"`
}{
    UserLocation: struct {
        City     string  `json:"city"`
        Distance float64 `json:"distance"`
    }{
        City:     "Palo Alto",
        Distance: 0.5,
    },
    PickUpToday: []Store{
        {
            ID:          "1",
            Title:       "Homeskillet Redwood City",
            Description: "Surprise Bag",
            PickUpTime:  "Pick up today 1:00 PM - 2:00 PM",
            Distance:    "3.8 mi",
            Price:       5.99,
            ImageURL:    "https://example.com/image1.jpg",
            Rating:      4.5,
            IsSaved:     false,
        },
        {
            ID:          "2",
            Title:       "Starbucks - University Ave",
            Description: "Surprise Bag",
            PickUpTime:  "Pick up today 5:00 PM - 6:00 PM",
            Distance:    "0.5 mi",
            Price:       4.99,
            ImageURL:    "https://example.com/image2.jpg",
            Rating:      4.2,
            IsSaved:     false,
        },
    },
    PickUpTomorrow: []Store{
        {
            ID:          "3",
            Title:       "Philz Coffee - Forest Ave",
            Description: "Surprise Bag",
            PickUpTime:  "Pick up tomorrow 7:00 AM - 8:00 AM",
            Distance:    "1.1 mi",
            Price:       3.99,
            ImageURL:    "https://example.com/image3.jpg",
            Rating:      4.3,
            IsSaved:     false,
        },
        {
            ID:          "4",
            Title:       "Blue Bottle Coffee",
            Description: "Surprise Bag",
            PickUpTime:  "Pick up tomorrow 8:00 AM - 9:00 AM",
            Distance:    "2.3 mi",
            Price:       6.99,
            ImageURL:    "https://example.com/image4.jpg",
            Rating:      4.7,
            IsSaved:     false,
        },
    },
    EmailVerified: true,
}

func SaveStore(c *gin.Context) {
    // TODO: In production, get actual userID from authentication
    userID := "test-user"
    storeID := c.Param("id")

    if savedStores[userID] == nil {
        savedStores[userID] = []string{}
    }

    // Check if store is already saved
    for _, id := range savedStores[userID] {
        if id == storeID {
            c.JSON(http.StatusOK, gin.H{"message": "Store already saved"})
            return
        }
    }

    savedStores[userID] = append(savedStores[userID], storeID)
    c.JSON(http.StatusOK, gin.H{"message": "Store saved successfully"})
}

func UnsaveStore(c *gin.Context) {
    userID := "test-user"
    storeID := c.Param("id")

    if savedStores[userID] == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No saved stores found"})
        return
    }

    // Remove store from saved stores
    for i, id := range savedStores[userID] {
        if id == storeID {
            savedStores[userID] = append(savedStores[userID][:i], savedStores[userID][i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "Store unsaved successfully"})
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "Store not found in saved stores"})
}

func GetFavorites(c *gin.Context) {
    userID := "test-user"
    userSavedStores := savedStores[userID]
    if userSavedStores == nil {
        c.JSON(http.StatusOK, []Store{})
        return
    }

    favorites := []Store{}
    for _, storeID := range userSavedStores {
        if store := findStoreByID(storeID); store != nil {
            favorites = append(favorites, *store)
        }
    }

    c.JSON(http.StatusOK, favorites)
}

func findStoreByID(storeID string) *Store {
    // Search in mock data (in production, this would be a database query)
    for _, store := range mockHomeResponse.PickUpToday {
        if store.ID == storeID {
            return &store
        }
    }
    for _, store := range mockHomeResponse.PickUpTomorrow {
        if store.ID == storeID {
            return &store
        }
    }
    return nil
} 