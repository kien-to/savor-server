package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"savor-server/db"
	"savor-server/models"

	"github.com/gin-gonic/gin"
)

type Store struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PickUpTime  string  `json:"pickUpTime"`
	Distance    string  `json:"distance"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
	Rating      float64 `json:"rating"`
	IsSaved     bool    `json:"isSaved"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type HomePageResponse struct {
	UserLocation struct {
		City     string `json:"city"`
		Distance int    `json:"distance"`
	} `json:"userLocation"`
	RecommendedStores []Store `json:"recommendedStores"`
	PickUpTomorrow    []Store `json:"pickUpTomorrow"`
	EmailVerified     bool    `json:"emailVerified"`
}

// @Summary     Get home page data
// @Description Get personalized home page data including recommended stores and pickup times
// @Tags        home
// @Accept      json
// @Produce     json
// @Param       latitude query number true "User's latitude"
// @Param       longitude query number true "User's longitude"
// @Success     200 {object} HomePageResponse
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Router      /api/home [get]
func GetHomePageData(c *gin.Context) {
	fmt.Println("GetHomePageData called")
	if db.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	lat := c.Query("latitude")
	lng := c.Query("longitude")

	if lat == "" || lng == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location parameters required"})
		return
	}

	var stores []models.Store
	err := db.DB.Select(&stores, `
		SELECT 
			id, title, description, pickup_time, 
			COALESCE(distance, '0 km') as distance,
			COALESCE(price, 0.0) as price, 
			COALESCE(original_price, 0.0) as original_price, 
			background_url, image_url,
			COALESCE(rating, 0.0) as rating, 
			COALESCE(reviews, 0) as reviews, 
			address, 
			COALESCE(items_left, 0) as items_left,
			latitude, longitude, created_at, updated_at,
			false as is_saved
		FROM stores 
		ORDER BY rating DESC 
		LIMIT 20
	`)

	if err != nil {
		fmt.Println(err)
		log.Printf("Failed to search stores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stores"})
		return
	}

	// Split stores into recommended and pickup tomorrow based on pickup time
	var recommended, tomorrow []models.Store
	for _, store := range stores {
		if store.PickupTime.Valid && strings.Contains(strings.ToLower(store.PickupTime.String), "tomorrow") {
			tomorrow = append(tomorrow, store)
		} else {
			recommended = append(recommended, store)
		}
	}

	response := HomePageResponse{
		UserLocation: struct {
			City     string `json:"city"`
			Distance int    `json:"distance"`
		}{
			City:     "Current Location",
			Distance: 5,
		},
		RecommendedStores: convertToStores(recommended),
		PickUpTomorrow:    convertToStores(tomorrow),
		EmailVerified:     true,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary     Search stores
// @Description Search for stores by name or description
// @Tags        home
// @Accept      json
// @Produce     json
// @Param       query query string true "Search query"
// @Success     200 {array} Store
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Router      /api/home/search [get]
func SearchStores(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	userID := c.GetString("user_id")
	var modelStores []models.Store
	err := db.DB.Select(&modelStores, `
		WITH saved_status AS (
			SELECT store_id, true as is_saved 
			FROM saved_stores 
			WHERE user_id = $1
		)
		SELECT s.*, 
			   array_agg(DISTINCT sh.highlight) FILTER (WHERE sh.highlight IS NOT NULL) as highlights,
			   COALESCE(ss.is_saved, false) as is_saved
		FROM stores s
		LEFT JOIN store_highlights sh ON s.id = sh.store_id
		LEFT JOIN saved_status ss ON s.id = ss.store_id
		WHERE 
			s.title ILIKE $2 OR 
			s.description ILIKE $2 OR 
			s.address ILIKE $2
		GROUP BY s.id, ss.is_saved
		ORDER BY s.rating DESC
		LIMIT 20
	`, userID, "%"+query+"%")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search stores"})
		return
	}

	// Convert model stores to response stores
	stores := convertToStores(modelStores)
	c.JSON(http.StatusOK, stores)
}

func convertToStores(modelStores []models.Store) []Store {
	stores := make([]Store, len(modelStores))
	for i, s := range modelStores {
		distance := "0 km"
		if s.Distance != nil {
			distance = *s.Distance
		}

		description := ""
		if s.Description.Valid {
			description = s.Description.String
		}

		pickupTime := ""
		if s.PickupTime.Valid {
			pickupTime = s.PickupTime.String
		}

		price := 0.0
		if s.Price.Valid {
			price = s.Price.Float64
		}

		rating := 0.0
		if s.Rating.Valid {
			rating = s.Rating.Float64
		}

		stores[i] = Store{
			ID:          s.ID,
			Title:       s.Title,
			Description: description,
			PickUpTime:  pickupTime,
			Distance:    distance,
			Price:       price,
			ImageURL:    s.ImageURL,
			Rating:      rating,
			IsSaved:     s.IsSaved,
			Latitude:    s.Latitude,
			Longitude:   s.Longitude,
		}
	}
	return stores
}
