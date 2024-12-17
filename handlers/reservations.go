package handlers

import (
	"fmt"
	"net/http"
	"savor-server/db"
	"time"

	"github.com/gin-gonic/gin"
)

type ReservationResponse struct {
	ID          string    `db:"id" json:"id"`
	StoreID     string    `db:"store_id" json:"storeId"`
	StoreName   string    `db:"store_name" json:"storeName"`
	StoreImage  string    `db:"store_image" json:"storeImage"`
	Quantity    int       `db:"quantity" json:"quantity"`
	TotalAmount float64   `db:"total_amount" json:"totalAmount"`
	Status      string    `db:"status" json:"status"`
	PaymentID   string    `db:"payment_id" json:"paymentId"`
	PickupTime  *string   `db:"pickup_time" json:"pickupTime,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
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
			r.created_at
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