package handlers

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"savor-server/db"
	"savor-server/services"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	// Register types with gob
	gob.Register([]ReservationResponse{})
	gob.Register(ReservationResponse{})
}

type ReservationResponse struct {
	ID              string     `db:"id" json:"id"`
	StoreID         string     `db:"store_id" json:"storeId"`
	StoreName       string     `db:"store_name" json:"storeName"`
	StoreImage      string     `db:"store_image" json:"storeImage"`
	StoreAddress    string     `db:"store_address" json:"storeAddress"`
	StoreLatitude   float64    `db:"store_latitude" json:"storeLatitude"`
	StoreLongitude  float64    `db:"store_longitude" json:"storeLongitude"`
	Quantity        int        `db:"quantity" json:"quantity"`
	TotalAmount     float64    `db:"total_amount" json:"totalAmount"`
	Status          string     `db:"status" json:"status"`
	PaymentID       string     `db:"payment_id" json:"paymentId"`
	PickupTime      *string    `db:"pickup_time" json:"pickupTime,omitempty"`
	PickupTimestamp *time.Time `db:"pickup_timestamp" json:"pickupTimestamp,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	OriginalPrice   float64    `db:"original_price" json:"originalPrice"`
	DiscountedPrice float64    `db:"discounted_price" json:"discountedPrice"`
	CustomerName    string     `db:"customer_name" json:"customerName,omitempty"`
	CustomerEmail   string     `db:"customer_email" json:"customerEmail,omitempty"`
	PhoneNumber     string     `db:"phone_number" json:"phoneNumber,omitempty"`
}

func GetUserReservations(c *gin.Context) {
	// Debug: Log all headers and context values
	log.Printf("DEBUG: GetUserReservations called")
	log.Printf("DEBUG: Authorization header: %s", c.GetHeader("Authorization"))
	log.Printf("DEBUG: Content-Type header: %s", c.GetHeader("Content-Type"))
	log.Printf("DEBUG: Request method: %s", c.Request.Method)
	log.Printf("DEBUG: Request URL: %s", c.Request.URL.String())

	// Check all context keys
	for key, value := range c.Keys {
		log.Printf("DEBUG: Context key '%s' = %v", key, value)
	}

	userID := c.GetString("user_id")
	log.Printf("DEBUG: Retrieved userID from context: '%s'", userID)

	if userID == "" {
		log.Printf("ERROR: User not authenticated - userID is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("DEBUG: GetUserReservations called for userID: %s", userID)

	// Check if reservations table exists
	var tableExists bool
	err := db.DB.Get(&tableExists, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'reservations'
		)
	`)

	if err != nil {
		log.Printf("ERROR: Failed to check if reservations table exists: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("DEBUG: Reservations table exists: %v", tableExists)

	if !tableExists {
		log.Printf("WARNING: Reservations table does not exist, returning empty array")
		// Return empty array since table doesn't exist yet
		c.JSON(http.StatusOK, []ReservationResponse{})
		return
	}

	// Query all reservations from database
	var allReservations []ReservationResponse
	err = db.DB.Select(&allReservations, `
		SELECT 
			r.id,
			r.store_id,
			s.title as store_name,
			s.image_url as store_image,
			s.address as store_address,
			s.latitude as store_latitude,
			s.longitude as store_longitude,
			r.quantity,
			r.total_amount,
			r.status,
			r.payment_id,
			r.pickup_time,
			r.pickup_timestamp,
			r.created_at,
			s.original_price,
			s.discounted_price,
			COALESCE(r.customer_name, u.email, 'Guest User') as customer_name,
			COALESCE(r.customer_email, u.email, '') as customer_email,
			COALESCE(r.phone_number, '') as phone_number
		FROM reservations r
		JOIN stores s ON r.store_id = s.id
		LEFT JOIN users u ON r.user_id = u.id::text
		WHERE r.user_id = $1 
		ORDER BY r.pickup_timestamp DESC, r.created_at DESC
	`, userID)

	if err != nil {
		log.Printf("ERROR: Failed to fetch reservations for userID %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	// Separate current and past reservations based on 24-hour window
	// Initialize as empty slices instead of nil to ensure JSON serialization as [] not null
	currentReservations := make([]ReservationResponse, 0)
	pastReservations := make([]ReservationResponse, 0)
	now := time.Now()
	twentyFourHoursAgo := now.Add(-24 * time.Hour)

	for _, reservation := range allReservations {
		if reservation.CreatedAt.After(twentyFourHoursAgo) {
			// Created within last 24 hours - current reservation
			currentReservations = append(currentReservations, reservation)
		} else {
			// Created more than 24 hours ago - past reservation
			pastReservations = append(pastReservations, reservation)
		}
	}

	log.Printf("DEBUG: Found %d current and %d past reservations for userID: %s", len(currentReservations), len(pastReservations), userID)

	c.JSON(http.StatusOK, gin.H{
		"currentReservations": currentReservations,
		"pastReservations":    pastReservations,
		"currentCount":        len(currentReservations),
		"pastCount":           len(pastReservations),
	})
}

func GetReservations(c *gin.Context) {
	// Check if user is authenticated by looking for Authorization header
	authHeader := c.GetHeader("Authorization")

	if authHeader != "" {
		// User has authorization header, try to get authenticated reservations
		// We need to manually verify the token here since we're not using middleware
		fmt.Println("Authorization header found, attempting authenticated request")

		// For now, return empty array for authenticated users since we don't have proper auth setup
		// In a real implementation, you would verify the token and get user reservations
		c.JSON(http.StatusOK, []ReservationResponse{})
		return
	}

	// No authorization header, treat as guest user and get session reservations
	fmt.Println("No authorization header, treating as guest user")
	GetGuestReservations(c)
}

func GetDemoReservations(c *gin.Context) {
	demoReservations := []ReservationResponse{
		{
			ID:              "demo-1",
			StoreID:         "store-1",
			StoreName:       "Sushi Paradise",
			StoreImage:      "https://images.unsplash.com/photo-1579871494447-9811cf80d66c",
			StoreAddress:    "123 Phan Chu Trinh, Hoàn Kiếm, Hà Nội",
			StoreLatitude:   21.0287,
			StoreLongitude:  105.8514,
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
			StoreAddress:    "456 Nguyen Hue, Hoan Kiem, Ha Noi",
			StoreLatitude:   21.0088,
			StoreLongitude:  105.8619,
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
}

func stringPtr(s string) *string {
	return &s
}

type GuestReservationRequest struct {
	StoreID         string  `json:"storeId"`
	StoreName       string  `json:"storeName"`
	StoreImage      string  `json:"storeImage"`
	StoreAddress    string  `json:"storeAddress"`
	StoreLatitude   float64 `json:"storeLatitude"`
	StoreLongitude  float64 `json:"storeLongitude"`
	Quantity        int     `json:"quantity"`
	TotalAmount     float64 `json:"totalAmount"`
	OriginalPrice   float64 `json:"originalPrice"`
	DiscountedPrice float64 `json:"discountedPrice"`
	PickupTime      string  `json:"pickupTime"`
	Name            string  `json:"name"`
	Email           string  `json:"email,omitempty"`
	Phone           string  `json:"phone,omitempty"`
	PaymentType     string  `json:"paymentType"`
}

func CreateAuthenticatedReservation(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		fmt.Printf("User not authenticated - userID is empty")
		log.Printf("ERROR: User not authenticated - userID is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req GuestReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Failed to bind JSON: %v", err)
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate required fields
	if req.StoreID == "" || req.Quantity < 1 {
		fmt.Printf("Invalid request data: %v", req)
		log.Printf("Invalid request data: %v", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	log.Printf("Creating authenticated reservation for user %s: %v", userID, req)
	fmt.Printf("Creating authenticated reservation for user %s: %v", userID, req)

	// Create a new reservation (use UUID for DB uuid type)
	reservation := ReservationResponse{
		ID:              uuid.New().String(),
		StoreID:         req.StoreID,
		StoreName:       req.StoreName,
		StoreImage:      req.StoreImage,
		StoreAddress:    req.StoreAddress,
		StoreLatitude:   req.StoreLatitude,
		StoreLongitude:  req.StoreLongitude,
		Quantity:        req.Quantity,
		TotalAmount:     req.TotalAmount,
		OriginalPrice:   req.OriginalPrice,
		DiscountedPrice: req.DiscountedPrice,
		Status:          "confirmed",
		PaymentID:       fmt.Sprintf("pay-%d", time.Now().Unix()),
		PickupTime:      &req.PickupTime,
		CreatedAt:       time.Now(),
		CustomerName:    req.Name,
		CustomerEmail:   req.Email,
		PhoneNumber:     req.Phone,
	}

	// Check if reservations table exists
	var tableExists bool
	err := db.DB.Get(&tableExists, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'reservations'
		)
	`)

	if err != nil {
		log.Printf("ERROR: Failed to check if reservations table exists: %v", err)
		fmt.Printf("ERROR: Failed to check if reservations table exists: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !tableExists {
		log.Printf("ERROR: Reservations table does not exist")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reservations table missing"})
		return
	}

	// Get store pickup timestamp directly from database
	var pickupTimestamp time.Time
	err = db.DB.Get(&pickupTimestamp, `SELECT pickup_timestamp FROM stores WHERE id = $1`, req.StoreID)
	if err != nil {
		log.Printf("WARNING: Failed to get store pickup timestamp for store %s: %v", req.StoreID, err)
		// Fallback: current time + 2 hours
		pickupTimestamp = time.Now().Add(2 * time.Hour)
	}

	// Insert into database
	_, err = db.DB.Exec(`
		INSERT INTO reservations (
			id, user_id, store_id, quantity, total_amount, 
			status, payment_id, pickup_time, pickup_timestamp, created_at,
			customer_name, customer_email, phone_number
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, reservation.ID, userID, req.StoreID, req.Quantity, req.TotalAmount,
		reservation.Status, reservation.PaymentID, req.PickupTime, pickupTimestamp, reservation.CreatedAt,
		req.Name, req.Email, req.Phone)

	if err != nil {
		log.Printf("ERROR: Failed to insert reservation into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	// Update items_left and bags_available count in stores table
	_, err = db.DB.Exec(`
		UPDATE stores 
		SET items_left = GREATEST(0, items_left - $1), 
		    bags_available = GREATEST(0, bags_available - $1), 
		    updated_at = NOW()
		WHERE id = $2
	`, req.Quantity, req.StoreID)

	if err != nil {
		log.Printf("WARNING: Failed to update items_left and bags_available for store %s: %v", req.StoreID, err)
		// Don't fail the reservation creation if items_left update fails
	} else {
		log.Printf("Updated items_left and bags_available for store %s: decreased by %d", req.StoreID, req.Quantity)
	}

	log.Printf("Reservation created successfully in database for user %s", userID)

	// Send notification (don't fail if notification fails)
	go func() {
		if services.NotificationSvc != nil {
			notificationData := services.ReservationNotificationData{
				CustomerName:  req.Name,
				StoreName:     req.StoreName,
				StoreAddress:  req.StoreAddress,
				Quantity:      req.Quantity,
				TotalAmount:   req.TotalAmount,
				PickupTime:    req.PickupTime,
				ReservationID: reservation.ID,
				Email:         req.Email,
				Phone:         req.Phone,
			}

			if err := services.NotificationSvc.SendReservationConfirmation(notificationData); err != nil {
				log.Printf("Failed to send notification: %v", err)
			} else {
				log.Printf("Notification sent successfully for reservation %s", reservation.ID)
			}
		}
	}()

	c.JSON(http.StatusOK, reservation)
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

	fmt.Printf("Creating guest reservation: %v\n", req)
	log.Printf("Creating guest reservation for store %s: %v", req.StoreID, req)

	// Get store pickup timestamp directly from database
	var pickupTimestamp time.Time
	err := db.DB.Get(&pickupTimestamp, `SELECT pickup_timestamp FROM stores WHERE id = $1`, req.StoreID)
	if err != nil {
		log.Printf("WARNING: Failed to get store pickup timestamp for store %s: %v", req.StoreID, err)
		// Fallback: current time + 2 hours
		pickupTimestamp = time.Now().Add(2 * time.Hour)
	}

	// Create a new reservation with UUID
	reservationID := uuid.New().String()
	reservation := ReservationResponse{
		ID:              reservationID,
		StoreID:         req.StoreID,
		StoreName:       req.StoreName,
		StoreImage:      req.StoreImage,
		StoreAddress:    req.StoreAddress,
		StoreLatitude:   req.StoreLatitude,
		StoreLongitude:  req.StoreLongitude,
		Quantity:        req.Quantity,
		TotalAmount:     req.TotalAmount,
		OriginalPrice:   req.OriginalPrice,
		DiscountedPrice: req.DiscountedPrice,
		Status:          "confirmed",
		PaymentID:       fmt.Sprintf("guest-pay-%d", time.Now().Unix()),
		PickupTime:      &req.PickupTime,
		PickupTimestamp: &pickupTimestamp,
		CreatedAt:       time.Now(),
		CustomerName:    req.Name,
		CustomerEmail:   req.Email,
		PhoneNumber:     req.Phone,
	}

	// Insert into database with NULL user_id for guest reservations
	_, err = db.DB.Exec(`
		INSERT INTO reservations (
			id, user_id, store_id, quantity, total_amount, 
			status, payment_id, pickup_time, pickup_timestamp, created_at,
			customer_name, customer_email, phone_number
		) VALUES ($1, NULL, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, reservationID, req.StoreID, req.Quantity, req.TotalAmount,
		reservation.Status, reservation.PaymentID, req.PickupTime, pickupTimestamp, reservation.CreatedAt,
		req.Name, req.Email, req.Phone)

	if err != nil {
		log.Printf("ERROR: Failed to insert guest reservation into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	// Update items_left and bags_available count in stores table
	_, err = db.DB.Exec(`
		UPDATE stores 
		SET items_left = GREATEST(0, items_left - $1), 
		    bags_available = GREATEST(0, bags_available - $1), 
		    updated_at = NOW()
		WHERE id = $2
	`, req.Quantity, req.StoreID)

	if err != nil {
		log.Printf("WARNING: Failed to update items_left and bags_available for store %s: %v", req.StoreID, err)
		// Don't fail the reservation creation if items_left update fails
	} else {
		log.Printf("Updated items_left and bags_available for store %s: decreased by %d", req.StoreID, req.Quantity)
	}

	log.Printf("Guest reservation created successfully in database: %s", reservationID)

	// Also store in session for backward compatibility
	session := sessions.Default(c)
	var sessionReservations []ReservationResponse
	if existingReservations := session.Get("reservations"); existingReservations != nil {
		sessionReservations = existingReservations.([]ReservationResponse)
	}
	sessionReservations = append(sessionReservations, reservation)
	session.Set("reservations", sessionReservations)
	if err := session.Save(); err != nil {
		log.Printf("WARNING: Failed to save guest reservation to session: %v", err)
		// Don't fail since it's already in the database
	}

	// Send notification (don't fail if notification fails)
	go func() {
		if services.NotificationSvc != nil {
			notificationData := services.ReservationNotificationData{
				CustomerName:  req.Name,
				StoreName:     req.StoreName,
				StoreAddress:  req.StoreAddress,
				Quantity:      req.Quantity,
				TotalAmount:   req.TotalAmount,
				PickupTime:    req.PickupTime,
				ReservationID: reservation.ID,
				Email:         req.Email,
				Phone:         req.Phone,
			}

			if err := services.NotificationSvc.SendReservationConfirmation(notificationData); err != nil {
				log.Printf("Failed to send notification: %v", err)
			} else {
				log.Printf("Notification sent successfully for reservation %s", reservation.ID)
			}
		}
	}()

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

	// For logged-in users, delete from database
	fmt.Printf("Attempting to delete reservation %s for user %s", reservationID, userID)
	log.Printf("Attempting to delete reservation %s for user %s", reservationID, userID)

	// First, get the reservation details to update bags_available
	var storeID string
	var quantity int
	err := db.DB.QueryRow(`
		SELECT store_id, quantity FROM reservations 
		WHERE id = $1 AND user_id = $2
	`, reservationID, userID).Scan(&storeID, &quantity)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Reservation %s not found for user %s", reservationID, userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
			return
		}
		fmt.Printf("ERROR: Failed to get reservation details %s: %v", reservationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reservation details"})
		return
	}

	// Delete the reservation
	result, err := db.DB.Exec(`
		DELETE FROM reservations 
		WHERE id = $1 AND user_id = $2
	`, reservationID, userID)

	if err != nil {
		fmt.Printf("ERROR: Failed to delete reservation %s: %v", reservationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reservation"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("ERROR: Failed to confirm deletion: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}

	if rowsAffected == 0 {
		fmt.Printf("Reservation %s not found for user %s", reservationID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	// Update items_left count in stores table (increase by the quantity that was reserved)
	_, err = db.DB.Exec(`
		UPDATE stores 
		SET items_left = items_left + $1, updated_at = NOW()
		WHERE id = $2
	`, quantity, storeID)

	if err != nil {
		log.Printf("WARNING: Failed to update items_left for store %s: %v", storeID, err)
		// Don't fail the deletion if items_left update fails
	} else {
		log.Printf("Updated items_left for store %s: increased by %d", storeID, quantity)
	}

	fmt.Printf("Successfully deleted reservation %s for user %s", reservationID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted"})
}

func DeleteGuestReservation(c *gin.Context) {
	reservationID := c.Param("id")

	log.Printf("Attempting to delete guest reservation %s", reservationID)

	// First, get the reservation details from database to update bags_available
	var storeID string
	var quantity int
	err := db.DB.QueryRow(`
		SELECT store_id, quantity FROM reservations 
		WHERE id = $1 AND user_id IS NULL
	`, reservationID).Scan(&storeID, &quantity)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Guest reservation %s not found in database", reservationID)
			// Try to delete from session anyway as fallback
			deleteFromSession(c, reservationID)
			return
		}
		log.Printf("ERROR: Failed to get guest reservation details %s: %v", reservationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reservation details"})
		return
	}

	// Delete from database
	result, err := db.DB.Exec(`
		DELETE FROM reservations 
		WHERE id = $1 AND user_id IS NULL
	`, reservationID)

	if err != nil {
		log.Printf("ERROR: Failed to delete guest reservation %s: %v", reservationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reservation"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Guest reservation %s not found", reservationID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	// Update items_left count in stores table (increase by the quantity that was reserved)
	_, err = db.DB.Exec(`
		UPDATE stores 
		SET items_left = items_left + $1, updated_at = NOW()
		WHERE id = $2
	`, quantity, storeID)

	if err != nil {
		log.Printf("WARNING: Failed to update items_left for store %s: %v", storeID, err)
		// Don't fail the deletion if items_left update fails
	} else {
		log.Printf("Updated items_left for store %s: increased by %d", storeID, quantity)
	}

	// Also delete from session for backward compatibility
	deleteFromSession(c, reservationID)

	log.Printf("Successfully deleted guest reservation %s", reservationID)
	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted"})
}

// Helper function to delete reservation from session
func deleteFromSession(c *gin.Context, reservationID string) {
	session := sessions.Default(c)
	if sessionReservations := session.Get("reservations"); sessionReservations != nil {
		reservations, ok := sessionReservations.([]ReservationResponse)
		if !ok {
			log.Printf("Failed to cast session reservations to []ReservationResponse: %v", sessionReservations)
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
		} else {
			log.Printf("Successfully deleted reservation %s from session", reservationID)
		}
	}
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

	fmt.Printf("Session reservations: %v\n", sessionReservations)

	reservations, ok := sessionReservations.([]ReservationResponse)
	if !ok {
		log.Printf("Failed to cast session reservations to []ReservationResponse: %v", sessionReservations)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
		return
	}

	// Debug each reservation's address
	for i, res := range reservations {
		fmt.Printf("[DEBUG] Reservation %d - ID: %s, StoreAddress: '%s'\n", i, res.ID, res.StoreAddress)
	}

	// Filter based on creation time (keep reservations created within last 24 hours)
	// This is more reliable than using pickup timestamp which might be outdated test data
	currentTime := time.Now()
	twentyFourHoursAgo := currentTime.Add(-24 * time.Hour)
	activeReservations := make([]ReservationResponse, 0)

	for _, res := range reservations {
		// Keep reservation if it was created within the last 24 hours
		if res.CreatedAt.After(twentyFourHoursAgo) {
			activeReservations = append(activeReservations, res)
			fmt.Printf("[DEBUG] Keeping reservation %s - created at %s\n", res.ID, res.CreatedAt)
		} else {
			fmt.Printf("[DEBUG] Filtering out reservation %s - created at %s (too old)\n", res.ID, res.CreatedAt)
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

// ClearSessionReservations clears all reservations from the session (for debugging)
func ClearSessionReservations(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("reservations")
	if err := session.Save(); err != nil {
		log.Printf("Error clearing session reservations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Session reservations cleared"})
}
