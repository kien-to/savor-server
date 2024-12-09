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
    BackgroundURL   string   `json:"backgroundUrl"`
    AvatarURL       string   `json:"avatarUrl"`
    Rating          float64  `json:"rating"`
    Reviews         int      `json:"reviews"`
    Address         string   `json:"address"`
    ItemsLeft       int      `json:"itemsLeft"`
    Highlights      []string `json:"highlights"`
    IsSaved         bool     `json:"isSaved"`
}

var mockStores = map[string]StoreDetail{
    "1": {
        ID:            "1",
        Title:         "Homeskillet Redwood City",
        Description:   "Surprise Bag",
        PickUpTime:    "Pick up tomorrow 1:00 AM - 5:00 AM",
        Distance:      "3.8 mi",
        Price:         5.99,
        OriginalPrice: 18.00,
        BackgroundURL: "https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png",
        AvatarURL:     "https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-logo.png",
        Rating:        4.6,
        Reviews:       952,
        Address:       "123 Main St, Redwood City, CA 94063",
        ItemsLeft:     5,
        Highlights:    []string{"Friendly staff", "Quick pickup", "Great value"},
    },
    "2": {
        ID:            "2",
        Title:         "Pho 75",
        Description:   "Surprise Bag",
        PickUpTime:    "Pick up tomorrow 1:00 AM - 5:00 AM",
        Distance:      "3.8 mi",
        Price:         5.99,
        OriginalPrice: 15.00,
        BackgroundURL: "https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg",
        AvatarURL:     "https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg",
        Rating:        4.6,
        Reviews:       1250,
        Address:       "456 Broadway St, Redwood City, CA 94063",
        ItemsLeft:     3,
        Highlights:    []string{"Authentic cuisine", "Large portions", "Fresh ingredients"},
    },
    "3": {
        ID:            "3",
        Title:         "Halal Guys",
        Description:   "Surprise Bag",
        PickUpTime:    "Pick up tomorrow 7:00 AM - 8:00 AM",
        Distance:      "1.1 mi",
        Price:         3.99,
        OriginalPrice: 12.00,
        BackgroundURL: "https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg",
        AvatarURL:     "https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg",
        Rating:        4.3,
        Reviews:       756,
        Address:       "789 El Camino Real, Menlo Park, CA 94025",
        ItemsLeft:     8,
        Highlights:    []string{"Halal certified", "Popular choice", "Fast service"},
    },
    "4": {
        ID:            "4",
        Title:         "Philz Coffee - Forest Ave",
        Description:   "Surprise Bag",
        PickUpTime:    "Pick up tomorrow 7:00 AM - 8:00 AM",
        Distance:      "1.1 mi",
        Price:         3.99,
        OriginalPrice: 10.00,
        BackgroundURL: "https://www.luxcafeclub.com/cdn/shop/articles/Minimalist_Modern_Coffee_Shop_1_1200x1200.png?v=1713243107",
        AvatarURL:     "https://www.luxcafeclub.com/cdn/shop/articles/Minimalist_Modern_Coffee_Shop_1_1200x1200.png?v=1713243107",
        Rating:        4.3,
        Reviews:       1543,
        Address:       "101 Forest Ave, Palo Alto, CA 94301",
        ItemsLeft:     4,
        Highlights:    []string{"Great coffee", "Cozy atmosphere", "Friendly baristas"},
    },
    "5": {
        ID:            "5",
        Title:         "Philz Coffee - Forest Ave",
        Description:   "Surprise Bag",
        PickUpTime:    "Pick up tomorrow 7:00 AM - 8:00 AM",
        Distance:      "1.1 mi",
        Price:         3.99,
        OriginalPrice: 10.00,
        BackgroundURL: "https://www.luxcafeclub.com/cdn/shop/articles/Minimalist_Modern_Coffee_Shop_1_1200x1200.png?v=1713243107",
        AvatarURL:     "https://www.luxcafeclub.com/cdn/shop/articles/Minimalist_Modern_Coffee_Shop_1_1200x1200.png?v=1713243107",
        Rating:        4.3,
        Reviews:       1543,
        Address:       "101 Forest Ave, Palo Alto, CA 94301",
        ItemsLeft:     4,
        Highlights:    []string{"Great coffee", "Cozy atmosphere", "Friendly baristas"},
    },
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
    
    store, exists := mockStores[storeId]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Store not found",
        })
        return
    }
    
    c.JSON(http.StatusOK, store)
}

// @Security    BearerAuth
// @Summary     Toggle store save status
// @Description Toggle whether a store is saved for the current user
// @Tags        stores
// @Accept      json
// @Produce     json
// @Param       id path string true "Store ID"
// @Success     200 {object} map[string]bool
// @Failure     404 {object} models.ErrorResponse
// @Router      /api/stores/{id}/toggle-save [post]
func ToggleSaveStore(c *gin.Context) {
    storeId := c.Param("id")
    
    store, exists := mockStores[storeId]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Store not found",
        })
        return
    }

    // Toggle the saved state
    store.IsSaved = !store.IsSaved
    mockStores[storeId] = store

    c.JSON(http.StatusOK, gin.H{
        "isSaved": store.IsSaved,
    })
} 