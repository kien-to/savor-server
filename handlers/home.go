package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
}

type HomePageResponse struct {
	UserLocation struct {
		City     string `json:"city"`
		Distance int    `json:"distance"` // in miles
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
	// Get user location from query params
	lat := c.Query("latitude")
	lng := c.Query("longitude")

	if lat == "" || lng == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location parameters required"})
		return
	}

	// TODO: Implement actual database queries
	// For now, return mock data that matches the frontend
	response := HomePageResponse{
		UserLocation: struct {
			City     string `json:"city"`
			Distance int    `json:"distance"`
		}{
			City:     "Menlo Park",
			Distance: 6,
		},
		RecommendedStores: []Store{
			{
				ID:          "1",
				Title:       "Homeskillet Redwood City",
				Description: "Surprise Bag",
				PickUpTime:  "Pick up tomorrow 1:00 AM - 5:00 AM",
				Distance:    "3.8 mi",
				Price:      5.99,
				ImageURL:   "/images/stores/homeskillet.jpg",
				Rating:     4.6,
			},
		},
		PickUpTomorrow: []Store{
			{
				ID:          "2",
				Title:       "Philz Coffee - Forest Ave",
				Description: "Surprise Bag",
				PickUpTime:  "Pick up tomorrow 7:00 AM - 8:00 AM",
				Distance:    "1.1 mi",
				Price:      3.99,
				ImageURL:   "/images/stores/philz.jpg",
				Rating:     4.3,
			},
		},
		EmailVerified: false,
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

	// TODO: Implement actual search logic
	// For now, return mock data
	stores := []Store{
		{
			ID:          "1",
			Title:       "Homeskillet Redwood City",
			Description: "Surprise Bag",
			PickUpTime:  "Pick up tomorrow 1:00 AM - 5:00 AM",
			Distance:    "3.8 mi",
			Price:      5.99,
			ImageURL:   "/images/stores/homeskillet.jpg",
			Rating:     4.6,
		},
	}

	c.JSON(http.StatusOK, stores)
} 