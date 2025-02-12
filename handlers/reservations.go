package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"savor-server/db"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	// Register types with gob
	gob.Register([]ReservationResponse{})
	gob.Register(ReservationResponse{})
}

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
		fmt.Printf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate required fields
	if req.StoreID == "" || req.Quantity < 1 || (req.Email == "" && req.Phone == "") {
		fmt.Printf("Invalid request data: %v\n", req)
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

	// Get session
	session := sessions.Default(c)
	var reservations []ReservationResponse

	// Get existing reservations from session
	if sessionReservations := session.Get("reservations"); sessionReservations != nil {
		reservations = sessionReservations.([]ReservationResponse)
	}

	// Add new reservation
	reservations = append(reservations, reservation)

	// Save to session
	session.Set("reservations", reservations)
	if err := session.Save(); err != nil {
		fmt.Printf("Failed to save reservation: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reservation"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func GetSessionReservations(c *gin.Context) {
	session := sessions.Default(c)
	if reservations := session.Get("reservations"); reservations != nil {
		c.JSON(http.StatusOK, reservations)
		return
	}
	c.JSON(http.StatusOK, []ReservationResponse{})
}

func DeleteReservation(c *gin.Context) {
	reservationID := c.Param("id")
	userID := c.GetString("user_id")

	// For guest users, check session
	if userID == "" {
		session := sessions.Default(c)
		if sessionReservations := session.Get("reservations"); sessionReservations != nil {
			reservations, ok := sessionReservations.([]ReservationResponse)
			if !ok {
				log.Printf("Failed to cast session reservations to []ReservationResponse: %v", sessionReservations)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
				return
			}

			newReservations := make([]ReservationResponse, 0)
			for _, r := range reservations {
				if r.ID != reservationID {
					newReservations = append(newReservations, r)
				}
			}

			session.Set("reservations", newReservations)

			if err := session.Save(); err != nil {
				log.Printf("Failed to save session after deleting reservation: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
				return
			}
		} else {
			log.Printf("No reservations found in session for deletion request: %s", reservationID)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted"})
		return
	}

	// For logged-in users, delete from database
	log.Printf("Attempting to delete reservation %s for user %s", reservationID, userID)

	result, err := db.DB.Exec(`
		DELETE FROM reservations 
		WHERE id = $1 AND user_id = $2
		AND (pickup_time IS NULL OR pickup_time > NOW())
	`, reservationID, userID)

	if err != nil {
		log.Printf("Database error deleting reservation %s: %v", reservationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reservation"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected for reservation %s: %v", reservationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}

	if rowsAffected == 0 {
		// Check if the reservation exists but is expired
		var exists bool
		err = db.DB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM reservations 
				WHERE id = $1 AND user_id = $2 AND pickup_time <= NOW()
			)
		`, reservationID, userID).Scan(&exists)

		if err != nil {
			log.Printf("Error checking if reservation %s is expired: %v", reservationID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check reservation status"})
			return
		}

		if exists {
			log.Printf("Attempted to delete expired reservation %s", reservationID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete expired reservation"})
			return
		}

		log.Printf("Reservation %s not found for user %s", reservationID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	log.Printf("Successfully deleted reservation %s for user %s", reservationID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted"})
}

// GetGuestReservations handles fetching reservations for guest users from their session
func GetGuestReservations(c *gin.Context) {
	session := sessions.Default(c)
	sessionReservations := session.Get("reservations")

	if sessionReservations == nil {
		// Return empty array instead of null
		c.JSON(http.StatusOK, []ReservationResponse{})
		return
	}

	reservations, ok := sessionReservations.([]ReservationResponse)
	if !ok {
		log.Printf("Failed to cast session reservations to []ReservationResponse: %v", sessionReservations)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
		return
	}

	// Filter out expired reservations
	currentTime := time.Now()
	activeReservations := make([]ReservationResponse, 0)
	for _, res := range reservations {
		if res.PickupTime == nil {
			activeReservations = append(activeReservations, res)
			continue
		}

		pickupTime, err := time.Parse(time.RFC3339, *res.PickupTime)
		if err != nil {
			log.Printf("Error parsing pickup time for reservation %s: %v", res.ID, err)
			continue
		}

		if pickupTime.After(currentTime) {
			activeReservations = append(activeReservations, res)
		}
	}

	// Update session with only active reservations
	if len(activeReservations) != len(reservations) {
		session.Set("reservations", activeReservations)
		if err := session.Save(); err != nil {
			log.Printf("Error saving session after filtering expired reservations: %v", err)
		}
	}

	c.JSON(http.StatusOK, activeReservations)
}
