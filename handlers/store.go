package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type StoreDetail struct {
    ID              string   `json:"id"`
    Title           string   `json:"title"`
    Description     string   `json:"description"`
    PickUpTime      string   `json:"pickUpTime"`
    Distance        string   `json:"distance"`
    Price           float64  `json:"price"`
    OriginalPrice   float64  `json:"originalPrice"`
    ImageURL        string   `json:"imageUrl"`
    Rating          float64  `json:"rating"`
    Reviews         int      `json:"reviews"`
    Address         string   `json:"address"`
    ItemsLeft       int      `json:"itemsLeft"`
    Highlights      []string `json:"highlights"`
}

// @Security    BearerAuth
// @Summary     Get store details
// @Description Get detailed information about a specific store
// @Tags        stores
// @Accept      json
// @Produce     json
// @Param       id path string true "Store ID"
// @Success     200 {object} StoreDetail
// @Failure     404 {object} models.ErrorResponse
// @Router      /api/stores/{id} [get]
func GetStoreDetail(c *gin.Context) {
    storeId := c.Param("id")
    
    // TODO: Implement actual database query
    // For now, return mock data
    store := StoreDetail{
        ID:            storeId,
        Title:         "Theo Chocolate",
        Description:   "Rescue a Surprise Bag filled with chocolates. Theo will also offer 10% off any regular priced purchases you make at the time of pick up.",
        PickUpTime:    "1:45 PM - 8:40 PM",
        Distance:      "0.5 mi",
        Price:         5.99,
        OriginalPrice: 18.00,
        ImageURL:      "https://example.com/store.jpg",
        Rating:        4.6,
        Reviews:       952,
        Address:       "3400 Phinney Ave N, Seattle, WA 98103, USA",
        ItemsLeft:     5,
        Highlights:    []string{"Friendly staff", "Quick pickup", "Great value"},
    }
    
    c.JSON(http.StatusOK, store)
} 