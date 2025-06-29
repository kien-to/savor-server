package handlers

import (
	"log"
	"net/http"
	"strconv"

	"savor-server/db"
	"savor-server/services"

	"github.com/gin-gonic/gin"
)

// @Summary     Calculate distance to store
// @Description Calculate the distance and duration from user's location to a store
// @Tags        maps
// @Accept      json
// @Produce     json
// @Param       userLat query number true "User's latitude"
// @Param       userLng query number true "User's longitude"
// @Param       storeLat query number true "Store's latitude"
// @Param       storeLng query number true "Store's longitude"
// @Success     200 {object} services.DistanceResult
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /api/maps/distance [get]
func CalculateDistance(c *gin.Context) {
	if services.GoogleMaps == nil {
		log.Printf("ERROR: CalculateDistance called but Google Maps service is not available.")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google Maps service not available"})
		return
	}

	userLatStr := c.Query("userLat")
	userLngStr := c.Query("userLng")
	storeLatStr := c.Query("storeLat")
	storeLngStr := c.Query("storeLng")

	log.Printf("Received /distance request: userLat=%s, userLng=%s, storeLat=%s, storeLng=%s", userLatStr, userLngStr, storeLatStr, storeLngStr)

	if userLatStr == "" || userLngStr == "" || storeLatStr == "" || storeLngStr == "" {
		log.Printf("ERROR: Missing location parameters in /distance request.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "All location parameters are required"})
		return
	}

	userLat, err := strconv.ParseFloat(userLatStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid user latitude '%s': %v", userLatStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user latitude"})
		return
	}

	userLng, err := strconv.ParseFloat(userLngStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid user longitude '%s': %v", userLngStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user longitude"})
		return
	}

	storeLat, err := strconv.ParseFloat(storeLatStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid store latitude '%s': %v", storeLatStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store latitude"})
		return
	}

	storeLng, err := strconv.ParseFloat(storeLngStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid store longitude '%s': %v", storeLngStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store longitude"})
		return
	}

	result, err := services.GoogleMaps.CalculateDistance(userLat, userLng, storeLat, storeLng)
	if err != nil {
		log.Printf("ERROR: service.CalculateDistance failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary     Get directions to store
// @Description Get directions and Google Maps URL from user's location to a store
// @Tags        maps
// @Accept      json
// @Produce     json
// @Param       userLat query number true "User's latitude"
// @Param       userLng query number true "User's longitude"
// @Param       storeLat query number true "Store's latitude"
// @Param       storeLng query number true "Store's longitude"
// @Success     200 {object} services.DirectionsResult
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /api/maps/directions [get]
func GetDirections(c *gin.Context) {
	if services.GoogleMaps == nil {
		log.Printf("ERROR: GetDirections called but Google Maps service is not available.")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google Maps service not available"})
		return
	}

	userLatStr := c.Query("userLat")
	userLngStr := c.Query("userLng")
	storeLatStr := c.Query("storeLat")
	storeLngStr := c.Query("storeLng")

	log.Printf("Received /directions request: userLat=%s, userLng=%s, storeLat=%s, storeLng=%s", userLatStr, userLngStr, storeLatStr, storeLngStr)

	if userLatStr == "" || userLngStr == "" || storeLatStr == "" || storeLngStr == "" {
		log.Printf("ERROR: Missing location parameters in /directions request.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "All location parameters are required"})
		return
	}

	userLat, err := strconv.ParseFloat(userLatStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid user latitude '%s': %v", userLatStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user latitude"})
		return
	}

	userLng, err := strconv.ParseFloat(userLngStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid user longitude '%s': %v", userLngStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user longitude"})
		return
	}

	storeLat, err := strconv.ParseFloat(storeLatStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid store latitude '%s': %v", storeLatStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store latitude"})
		return
	}

	storeLng, err := strconv.ParseFloat(storeLngStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid store longitude '%s': %v", storeLngStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store longitude"})
		return
	}

	result, err := services.GoogleMaps.GetDirections(userLat, userLng, storeLat, storeLng)
	if err != nil {
		log.Printf("ERROR: service.GetDirections failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary     Get store with distance and directions
// @Description Get store details with calculated distance and directions URL
// @Tags        maps
// @Accept      json
// @Produce     json
// @Param       storeId path string true "Store ID"
// @Param       userLat query number true "User's latitude"
// @Param       userLng query number true "User's longitude"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     404 {object} map[string]string "Store not found"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /api/maps/stores/{storeId} [get]
func GetStoreWithDistance(c *gin.Context) {
	storeID := c.Param("storeId")
	userLatStr := c.Query("userLat")
	userLngStr := c.Query("userLng")

	log.Printf("Received /stores/%s request: userLat=%s, userLng=%s", storeID, userLatStr, userLngStr)

	if storeID == "" {
		log.Printf("ERROR: Missing store ID in /stores/{storeId} request.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Store ID is required"})
		return
	}

	if userLatStr == "" || userLngStr == "" {
		log.Printf("ERROR: Missing user location in /stores/%s request.", storeID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User location is required"})
		return
	}

	userLat, err := strconv.ParseFloat(userLatStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid user latitude '%s' for store %s: %v", userLatStr, storeID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user latitude"})
		return
	}

	userLng, err := strconv.ParseFloat(userLngStr, 64)
	if err != nil {
		log.Printf("ERROR: Invalid user longitude '%s' for store %s: %v", userLngStr, storeID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user longitude"})
		return
	}

	// Get store details first
	var store struct {
		ID        string  `db:"id"`
		Title     string  `db:"title"`
		Address   string  `db:"address"`
		Latitude  float64 `db:"latitude"`
		Longitude float64 `db:"longitude"`
	}

	err = db.DB.Get(&store, `
		SELECT id, title, address, latitude, longitude
		FROM stores 
		WHERE id = $1
	`, storeID)

	if err != nil {
		log.Printf("ERROR: Could not find store with ID %s: %v", storeID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	// Calculate distance and get directions
	var distanceResult *services.DistanceResult
	var directionsResult *services.DirectionsResult

	if services.GoogleMaps != nil {
		distanceResult, err = services.GoogleMaps.CalculateDistance(userLat, userLng, store.Latitude, store.Longitude)
		if err != nil {
			log.Printf("ERROR: Could not calculate distance for store %s: %v", storeID, err)
			// Don't return, just log the error and continue
		}

		directionsResult, err = services.GoogleMaps.GetDirections(userLat, userLng, store.Latitude, store.Longitude)
		if err != nil {
			log.Printf("ERROR: Could not get directions for store %s: %v", storeID, err)
			// Don't return, just log the error and continue
		}
	} else {
		log.Printf("WARNING: Google Maps service is not available for store %s.", storeID)
	}

	response := gin.H{
		"store": store,
	}

	if distanceResult != nil {
		response["distance"] = distanceResult
	}

	if directionsResult != nil {
		response["directions"] = directionsResult
	}

	c.JSON(http.StatusOK, response)
}
