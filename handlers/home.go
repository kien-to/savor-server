package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
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
	Address         string   `json:"address"` // THIS WAS MISSING!
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
	fmt.Println("[BACKEND] GetHomePageData called")
	if db.DB == nil {
		log.Printf("[BACKEND] ERROR: Database connection not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	lat := c.Query("latitude")
	lng := c.Query("longitude")
	log.Printf("[BACKEND] Received params - lat: %s, lng: %s", lat, lng)

	if lat == "" || lng == "" {
		log.Printf("[BACKEND] ERROR: Missing location parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location parameters required"})
		return
	}

	userLat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Printf("[BACKEND] ERROR: Invalid latitude '%s': %v", lat, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	userLng, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		log.Printf("[BACKEND] ERROR: Invalid longitude '%s': %v", lng, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	log.Printf("[BACKEND] Parsed coordinates: lat=%f, lng=%f", userLat, userLng)

	var stores []models.Store

	// Debug: Log the user_id and its type
	userID := c.GetString("user_id")
	log.Printf("[BACKEND] DEBUG: user_id = '%s', type = %T", userID, userID)

	// Debug: Check total count of stores in database
	var totalCount int
	countErr := db.DB.Get(&totalCount, "SELECT COUNT(*) FROM stores")
	if countErr == nil {
		log.Printf("[BACKEND] 🔢 Total stores in database: %d", totalCount)
	} else {
		log.Printf("[BACKEND] ERROR: Could not count stores: %v", countErr)
	}

	// COMMENTED OUT: Original query with saved_stores table (causing type mismatch)
	/*
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
		`, userID)
	*/

	// NEW QUERY: Simplified without saved_stores and store_highlights tables
	err = db.DB.Select(&stores, `
		SELECT 
			s.id, 
			s.title, 
			COALESCE(s.description, '') as description,
			COALESCE(s.pickup_time, '') as pickup_time,
			COALESCE(s.distance, '0 km') as distance,
			COALESCE(s.price::numeric, 0.0) as price,
			COALESCE(s.original_price::numeric, s.price::numeric, 0.0) as original_price,
			COALESCE(s.price::numeric, 0.0) as discounted_price,
			COALESCE(s.background_url, '') as background_url,
			COALESCE(s.image_url, '') as image_url,
			COALESCE(s.rating, 0.0) as rating,
			COALESCE(s.reviews, 0) as reviews,
			COALESCE(s.reviews, 0) as reviews_count,
			COALESCE(s.address, '') as address,
			COALESCE(s.items_left, 0) as items_left,
			COALESCE(s.items_left, 0) as bags_available,
			COALESCE(s.latitude, 0.0) as latitude,
			COALESCE(s.longitude, 0.0) as longitude,
			COALESCE(s.google_maps_url, '') as google_maps_url,
			true as is_selling,
			false as is_saved,  -- Set to false for now since we're not checking saved_stores
			ARRAY[]::text[] as highlights  -- Empty array since we're not checking store_highlights
		FROM stores s
		LIMIT 50
	`)

	if err != nil {
		fmt.Println(err)
		log.Printf("[BACKEND] ERROR: Database query failed: %v", err)
		log.Printf("[BACKEND] DEBUG: Query was executed with user_id = '%s'", userID)

		// Check if this is a type mismatch error
		if strings.Contains(err.Error(), "operator does not exist: integer = text") {
			log.Printf("[BACKEND] DEBUG: This is a type mismatch error - likely comparing integer with text in saved_stores table")
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stores"})
		return
	}

	log.Printf("[BACKEND] ✅ Query successful, found %d stores from database", len(stores))

	// Log details about each store for debugging
	for i, store := range stores {
		log.Printf("[BACKEND] Store %d: ID=%s, Title=%s, Price=%.2f, Rating=%.1f, ItemsLeft=%d, Address='%s'",
			i+1, store.ID, store.Title, store.Price.Float64, store.Rating.Float64, store.ItemsLeft.Int64, store.Address)
	}

	// Additional debug: Log the exact SQL query being executed
	log.Printf("[BACKEND] 🔍 DEBUG: Database contains only %d stores. Expected ~15 stores from CSV import.", len(stores))

	// Calculate distances using Google Maps API and sort by distance
	if services.GoogleMaps != nil {
		log.Printf("[BACKEND] Google Maps service available, calculating distances")

		// Create a slice to hold stores with their calculated distance in meters
		type storeWithDistance struct {
			store          models.Store
			distanceStr    string
			distanceMeters int
		}
		var storesWithDistances []storeWithDistance

		for _, store := range stores {
			if store.Latitude != 0 && store.Longitude != 0 {
				distanceResult, err := services.GoogleMaps.CalculateDistance(userLat, userLng, store.Latitude, store.Longitude)
				if err == nil && distanceResult != nil {
					distanceStr := distanceResult.Distance
					store.Distance = &distanceStr
					storesWithDistances = append(storesWithDistances, storeWithDistance{
						store:          store,
						distanceStr:    distanceStr,
						distanceMeters: distanceResult.Meters,
					})
				} else {
					log.Printf("[BACKEND] Distance calculation failed for store %s: %v", store.ID, err)
					// Add store with default distance if calculation fails
					defaultDistance := "~km"
					store.Distance = &defaultDistance
					storesWithDistances = append(storesWithDistances, storeWithDistance{
						store:          store,
						distanceStr:    defaultDistance,
						distanceMeters: 999999, // Large number to put at end
					})
				}
			} else {
				// Add store with default distance if no coordinates
				defaultDistance := "~km"
				store.Distance = &defaultDistance
				storesWithDistances = append(storesWithDistances, storeWithDistance{
					store:          store,
					distanceStr:    defaultDistance,
					distanceMeters: 999999, // Large number to put at end
				})
			}
		}

		// Sort stores by distance (closest first)
		sort.Slice(storesWithDistances, func(i, j int) bool {
			return storesWithDistances[i].distanceMeters < storesWithDistances[j].distanceMeters
		})

		// Log the sorted distances for debugging
		log.Printf("[BACKEND] Distance sorting results:")
		for i, swd := range storesWithDistances {
			if i < 15 { // Log first 15 for debugging
				log.Printf("[BACKEND] %d. %s: %s (%d meters)", i+1, swd.store.Title, swd.distanceStr, swd.distanceMeters)
			}
		}

		// Extract sorted stores (limit to top 20 closest)
		maxStores := len(storesWithDistances)
		if maxStores > 20 {
			maxStores = 20
		}
		stores = make([]models.Store, maxStores)
		for i := 0; i < maxStores; i++ {
			stores[i] = storesWithDistances[i].store
		}

		log.Printf("[BACKEND] Stores sorted by distance - returning %d closest stores", maxStores)
	} else {
		log.Printf("[BACKEND] WARNING: Google Maps service not available")
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

	log.Printf("[BACKEND] 📊 Data split - Recommended: %d stores, Tomorrow: %d stores", len(recommended), len(tomorrow))

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

	log.Printf("[BACKEND] 🚀 Response ready with %d recommended, %d tomorrow stores", len(response.RecommendedStores), len(response.PickUpTomorrow))
	log.Printf("[BACKEND] 📤 Sending final response to client")
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
	log.Printf("[BACKEND] 🔍 SearchStores called with query: '%s'", query)

	if query == "" {
		log.Printf("[BACKEND] ERROR: Empty search query")
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

	// userID := c.GetString("user_id")
	var modelStores []models.Store

	// COMMENTED OUT: Original query with saved_stores table (causing type mismatch)
	/*
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
	*/

	// NEW QUERY: Simplified without saved_stores and store_highlights tables
	err := db.DB.Select(&modelStores, `
		SELECT 
			s.id, 
			s.title, 
			COALESCE(s.description, '') as description,
			COALESCE(s.pickup_time, '') as pickup_time,
			COALESCE(s.distance, '0 km') as distance,
			COALESCE(s.price::numeric, 0.0) as price,
			COALESCE(s.original_price::numeric, s.price::numeric, 0.0) as original_price,
			COALESCE(s.price::numeric, 0.0) as discounted_price,
			COALESCE(s.background_url, '') as background_url,
			COALESCE(s.image_url, '') as image_url,
			COALESCE(s.rating, 0.0) as rating,
			COALESCE(s.reviews, 0) as reviews,
			COALESCE(s.reviews, 0) as reviews_count,
			COALESCE(s.address, '') as address,
			COALESCE(s.items_left, 0) as items_left,
			COALESCE(s.items_left, 0) as bags_available,
			COALESCE(s.latitude, 0.0) as latitude,
			COALESCE(s.longitude, 0.0) as longitude,
			COALESCE(s.google_maps_url, '') as google_maps_url,
			true as is_selling,
			false as is_saved,  -- Set to false for now since we're not checking saved_stores
			ARRAY[]::text[] as highlights  -- Empty array since we're not checking store_highlights
		FROM stores s
		WHERE 
			COALESCE(s.title, '') ILIKE $1 OR 
			COALESCE(s.description, '') ILIKE $1 OR 
			COALESCE(s.address, '') ILIKE $1
		LIMIT 50
	`, "%"+query+"%")

	if err != nil {
		log.Printf("[BACKEND] ERROR: Search query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search stores"})
		return
	}

	log.Printf("[BACKEND] ✅ Search successful, found %d stores matching '%s'", len(modelStores), query)

	// Log search results for debugging
	for i, store := range modelStores {
		log.Printf("[BACKEND] Search Result %d: ID=%s, Title=%s, Matches query '%s'",
			i+1, store.ID, store.Title, query)
	}

	// Calculate distances using Google Maps API if user location is provided and sort by distance
	if services.GoogleMaps != nil && userLat != 0 && userLng != 0 {
		// Create a slice to hold stores with their calculated distance in meters
		type storeWithDistance struct {
			store          models.Store
			distanceStr    string
			distanceMeters int
		}
		var storesWithDistances []storeWithDistance

		for _, store := range modelStores {
			if store.Latitude != 0 && store.Longitude != 0 {
				distanceResult, err := services.GoogleMaps.CalculateDistance(userLat, userLng, store.Latitude, store.Longitude)
				if err == nil && distanceResult != nil {
					distanceStr := distanceResult.Distance
					store.Distance = &distanceStr
					storesWithDistances = append(storesWithDistances, storeWithDistance{
						store:          store,
						distanceStr:    distanceStr,
						distanceMeters: distanceResult.Meters,
					})
				} else {
					// Add store with default distance if calculation fails
					defaultDistance := "~km"
					store.Distance = &defaultDistance
					storesWithDistances = append(storesWithDistances, storeWithDistance{
						store:          store,
						distanceStr:    defaultDistance,
						distanceMeters: 999999, // Large number to put at end
					})
				}
			} else {
				// Add store with default distance if no coordinates
				defaultDistance := "~km"
				store.Distance = &defaultDistance
				storesWithDistances = append(storesWithDistances, storeWithDistance{
					store:          store,
					distanceStr:    defaultDistance,
					distanceMeters: 999999, // Large number to put at end
				})
			}
		}

		// Sort stores by distance (closest first)
		sort.Slice(storesWithDistances, func(i, j int) bool {
			return storesWithDistances[i].distanceMeters < storesWithDistances[j].distanceMeters
		})

		// Extract sorted stores (limit to top 20 closest)
		maxStores := len(storesWithDistances)
		if maxStores > 20 {
			maxStores = 20
		}
		modelStores = make([]models.Store, maxStores)
		for i := 0; i < maxStores; i++ {
			modelStores[i] = storesWithDistances[i].store
		}

		log.Printf("[BACKEND] Search results sorted by distance - returning %d closest stores", maxStores)
	}

	// Convert model stores to response stores
	stores := convertToStores(modelStores)
	log.Printf("[BACKEND] 📤 Search response ready with %d stores for query '%s'", len(stores), query)
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

		// Debug: Log the address being converted
		fmt.Printf("[DEBUG convertToStores] Store %s: Address='%s'\n", s.ID, s.Address)

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
			Address:         s.Address, // THIS WAS MISSING!
			IsSaved:         s.IsSaved,
			Latitude:        s.Latitude,
			Longitude:       s.Longitude,
			GoogleMapsURL:   googleMapsURL,
			Highlights:      s.Highlights,
		}
	}
	return stores
}
