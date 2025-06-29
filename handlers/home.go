package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"savor-server/db"
	"savor-server/models"
	"savor-server/services"

	"github.com/gin-gonic/gin"
)

type Store struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	PickUpTime      string   `json:"pickUpTime"`
	Distance        string   `json:"distance"`
	Price           float64  `json:"price"`
	OriginalPrice   float64  `json:"originalPrice"`
	DiscountedPrice float64  `json:"discountedPrice"`
	ImageURL        string   `json:"imageUrl"`
	Rating          float64  `json:"rating"`
	IsSaved         bool     `json:"isSaved"`
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
	GoogleMapsURL   string   `json:"googleMapsUrl"`
	ReviewsCount    int64    `json:"reviewsCount"`
	BagsAvailable   int64    `json:"bagsAvailable"`
	Highlights      []string `json:"highlights"`
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

	userLat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	userLng, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	var stores []models.Store
	err = db.DB.Select(&stores, `
		WITH saved_status AS (
			SELECT store_id, true as is_saved 
			FROM saved_stores 
			WHERE user_id = $1
		)
		SELECT 
			s.id, 
			s.title, 
			s.description, 
			s.pickup_time,
			COALESCE(s.distance, '0 km') as distance,
			COALESCE(s.price, 0.0) as price,
			COALESCE(s.original_price, s.price) as original_price,
			COALESCE(s.discounted_price, s.price) as discounted_price,
			s.background_url,
			s.image_url,
			COALESCE(s.rating, 0.0) as rating,
			COALESCE(s.reviews, 0) as reviews,
			COALESCE(s.reviews_count, s.reviews) as reviews_count,
			s.address,
			COALESCE(s.items_left, 0) as items_left,
			COALESCE(s.bags_available, s.items_left) as bags_available,
			s.latitude,
			s.longitude,
			s.google_maps_url,
			s.is_selling,
			COALESCE(ss.is_saved, false) as is_saved,
			array_agg(DISTINCT sh.highlight) FILTER (WHERE sh.highlight IS NOT NULL) as highlights
		FROM stores s
		LEFT JOIN saved_status ss ON s.id = ss.store_id
		LEFT JOIN store_highlights sh ON s.id = sh.store_id
		WHERE s.is_selling = true
		GROUP BY 
			s.id, 
			s.title,
			s.description,
			s.pickup_time,
			s.distance,
			s.price,
			s.original_price,
			s.discounted_price,
			s.background_url,
			s.image_url,
			s.rating,
			s.reviews,
			s.reviews_count,
			s.address,
			s.items_left,
			s.bags_available,
			s.latitude,
			s.longitude,
			s.google_maps_url,
			s.is_selling,
			ss.is_saved
		ORDER BY s.rating DESC 
		LIMIT 20
	`, c.GetString("user_id"))

	if err != nil {
		fmt.Println(err)
		log.Printf("Failed to search stores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stores"})
		return
	}

	// Calculate distances using Google Maps API
	if services.GoogleMaps != nil {
		for i := range stores {
			if stores[i].Latitude != 0 && stores[i].Longitude != 0 {
				distanceResult, err := services.GoogleMaps.CalculateDistance(userLat, userLng, stores[i].Latitude, stores[i].Longitude)
				if err == nil && distanceResult != nil {
					distanceStr := distanceResult.Distance
					stores[i].Distance = &distanceStr
				}
			}
		}
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

	lat := c.Query("latitude")
	lng := c.Query("longitude")

	userLat := 0.0
	userLng := 0.0
	if lat != "" && lng != "" {
		var err error
		userLat, err = strconv.ParseFloat(lat, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}

		userLng, err = strconv.ParseFloat(lng, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}
	}

	userID := c.GetString("user_id")
	var modelStores []models.Store
	err := db.DB.Select(&modelStores, `
		WITH saved_status AS (
			SELECT store_id, true as is_saved 
			FROM saved_stores 
			WHERE user_id = $1
		)
		SELECT 
			s.id, 
			s.title, 
			s.description, 
			s.pickup_time,
			COALESCE(s.distance, '0 km') as distance,
			COALESCE(s.price, 0.0) as price,
			COALESCE(s.original_price, s.price) as original_price,
			COALESCE(s.discounted_price, s.price) as discounted_price,
			s.background_url,
			s.image_url,
			COALESCE(s.rating, 0.0) as rating,
			COALESCE(s.reviews, 0) as reviews,
			COALESCE(s.reviews_count, s.reviews) as reviews_count,
			s.address,
			COALESCE(s.items_left, 0) as items_left,
			COALESCE(s.bags_available, s.items_left) as bags_available,
			s.latitude,
			s.longitude,
			s.google_maps_url,
			s.is_selling,
			COALESCE(ss.is_saved, false) as is_saved,
			array_agg(DISTINCT sh.highlight) FILTER (WHERE sh.highlight IS NOT NULL) as highlights
		FROM stores s
		LEFT JOIN saved_status ss ON s.id = ss.store_id
		LEFT JOIN store_highlights sh ON s.id = sh.store_id
		WHERE 
			s.title ILIKE $2 OR 
			s.description ILIKE $2 OR 
			s.address ILIKE $2
		GROUP BY 
			s.id, 
			s.title,
			s.description,
			s.pickup_time,
			s.distance,
			s.price,
			s.original_price,
			s.discounted_price,
			s.background_url,
			s.image_url,
			s.rating,
			s.reviews,
			s.reviews_count,
			s.address,
			s.items_left,
			s.bags_available,
			s.latitude,
			s.longitude,
			s.google_maps_url,
			s.is_selling,
			ss.is_saved
		ORDER BY s.rating DESC
		LIMIT 20
	`, userID, "%"+query+"%")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search stores"})
		return
	}

	// Calculate distances using Google Maps API if user location is provided
	if services.GoogleMaps != nil && userLat != 0 && userLng != 0 {
		for i := range modelStores {
			if modelStores[i].Latitude != 0 && modelStores[i].Longitude != 0 {
				distanceResult, err := services.GoogleMaps.CalculateDistance(userLat, userLng, modelStores[i].Latitude, modelStores[i].Longitude)
				if err == nil && distanceResult != nil {
					distanceStr := distanceResult.Distance
					modelStores[i].Distance = &distanceStr
				}
			}
		}
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

		originalPrice := price
		if s.OriginalPrice.Valid {
			originalPrice = s.OriginalPrice.Float64
		}

		discountedPrice := price
		if s.DiscountedPrice.Valid {
			discountedPrice = s.DiscountedPrice.Float64
		}

		rating := 0.0
		if s.Rating.Valid {
			rating = s.Rating.Float64
		}

		reviewsCount := int64(0)
		if s.ReviewsCount.Valid {
			reviewsCount = s.ReviewsCount.Int64
		}

		bagsAvailable := int64(0)
		if s.BagsAvailable.Valid {
			bagsAvailable = s.BagsAvailable.Int64
		}

		googleMapsURL := ""
		if s.GoogleMapsURL.Valid {
			googleMapsURL = s.GoogleMapsURL.String
		}

		stores[i] = Store{
			ID:              s.ID,
			Title:           s.Title,
			Description:     description,
			PickUpTime:      pickupTime,
			Distance:        distance,
			Price:           price,
			OriginalPrice:   originalPrice,
			DiscountedPrice: discountedPrice,
			ImageURL:        s.ImageURL,
			Rating:          rating,
			ReviewsCount:    reviewsCount,
			BagsAvailable:   bagsAvailable,
			IsSaved:         s.IsSaved,
			Latitude:        s.Latitude,
			Longitude:       s.Longitude,
			GoogleMapsURL:   googleMapsURL,
			Highlights:      s.Highlights,
		}
	}
	return stores
}
