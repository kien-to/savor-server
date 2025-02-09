package handlers

import (
	"fmt"
	"net/http"
	"savor-server/db"
	"time"

	"github.com/gin-gonic/gin"
)

type ReservationResponse struct {
	ID              string    `db:"id" json:"id"`
	StoreID         string    `db:"store_id" json:"storeId"`
	StoreName       string    `db:"store_name" json:"storeName"`
	StoreImage      string    `db:"store_image" json:"storeImage"`
	Quantity        int       `db:"quantity" json:"quantity"`
	TotalAmount     float64   `db:"total_amount" json:"totalAmount"`
	Status          string    `db:"status" json:"status"`
	PaymentID       string    `db:"payment_id" json:"paymentId"`
	PickupTime      *string   `db:"pickup_time" json:"pickupTime,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	OriginalPrice   float64   `db:"original_price" json:"originalPrice"`
	DiscountedPrice float64   `db:"discounted_price" json:"discountedPrice"`
}

func GetUserReservations(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		fmt.Println("User not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var reservations []ReservationResponse
	err := db.DB.Select(&reservations, `
		SELECT 
			r.id,
			r.store_id,
			s.title as store_name,
			s.image_url as store_image,
			r.quantity,
			r.total_amount,
			r.status,
			r.payment_id,
			r.pickup_time,
			r.created_at,
			r.original_price,
			r.discounted_price
		FROM reservations r
		JOIN stores s ON r.store_id = s.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
	`, userID)

	if err != nil {
		fmt.Printf("Failed to fetch reservations: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func GetDemoReservations(c *gin.Context) {
	demoReservations := []ReservationResponse{
		{
			ID:              "demo-1",
			StoreID:         "store-1",
			StoreName:       "Sushi Paradise",
			StoreImage:      "https://images.unsplash.com/photo-1579871494447-9811cf80d66c",
			Quantity:        2,
			TotalAmount:     19.99,
			Status:          "confirmed",
			PaymentID:       "pay-1",
			PickupTime:      stringPtr("2024-02-09T18:00:00Z"),
			CreatedAt:       time.Now().Add(-24 * time.Hour),
			OriginalPrice:   29.99,
			DiscountedPrice: 19.99,
		},
		{
			ID:              "demo-2",
			StoreID:         "store-2",
			StoreName:       "Bakery Delight",
			StoreImage:      "https://images.unsplash.com/photo-1509440159596-0249088772ff",
			Quantity:        1,
			TotalAmount:     9.99,
			Status:          "completed",
			PaymentID:       "pay-2",
			PickupTime:      stringPtr("2024-02-08T17:00:00Z"),
			CreatedAt:       time.Now().Add(-48 * time.Hour),
			OriginalPrice:   15.99,
			DiscountedPrice: 9.99,
		},
	}

	c.JSON(http.StatusOK, demoReservations)

	// print errors

}

func stringPtr(s string) *string {
	return &s
}

type GuestReservationRequest struct {
	StoreID     string  `json:"storeId"`
	StoreName   string  `json:"storeName"`
	StoreImage  string  `json:"storeImage"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"totalAmount"`
	PickupTime  string  `json:"pickupTime"`
	Name        string  `json:"name"`
	Email       string  `json:"email,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	PaymentType string  `json:"paymentType"`
}

func CreateGuestReservation(c *gin.Context) {
	var req GuestReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate required fields
	if req.StoreID == "" || req.Quantity < 1 || (req.Email == "" && req.Phone == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Create a new reservation
	reservation := ReservationResponse{
		ID:          fmt.Sprintf("guest-%d", time.Now().Unix()),
		StoreID:     req.StoreID,
		StoreName:   req.StoreName,
		StoreImage:  req.StoreImage,
		Quantity:    req.Quantity,
		TotalAmount: req.TotalAmount,
		Status:      "pending",
		PaymentID:   fmt.Sprintf("pay-%d", time.Now().Unix()),
		PickupTime:  &req.PickupTime,
		CreatedAt:   time.Now(),
	}

	// Get existing reservations from session or initialize new slice
	session := c.MustGet("session").(map[string]interface{})
	var reservations []ReservationResponse
	if existingReservations, ok := session["reservations"].([]ReservationResponse); ok {
		reservations = existingReservations
	}

	// Add new reservation
	reservations = append(reservations, reservation)
	session["reservations"] = reservations

	c.JSON(http.StatusOK, reservation)
}

func GetSessionReservations(c *gin.Context) {
	session := c.MustGet("session").(map[string]interface{})
	if reservations, ok := session["reservations"].([]ReservationResponse); ok {
		c.JSON(http.StatusOK, reservations)
		return
	}
	c.JSON(http.StatusOK, []ReservationResponse{})
}
