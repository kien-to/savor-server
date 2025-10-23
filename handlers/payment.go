package handlers

import (
	"fmt"
	"net/http"
	"savor-server/db"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type ReservationRequest struct {
	StoreId       string  `json:"storeId" binding:"required"`
	Quantity      int     `json:"quantity" binding:"required,min=1"`
	TotalAmount   float64 `json:"totalAmount" binding:"required"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
	PickupTime    string  `json:"pickupTime" binding:"required"`
}

type PayAtStoreRequest struct {
	PaymentIntentId string `json:"paymentIntentId" binding:"required"`
	PaymentMethod   string `json:"paymentMethod" binding:"required"`
}

func CreateReservation(c *gin.Context) {
	var req ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Request binding error: %+v\n", err)
		fmt.Printf("Received request body: %+v\n", req)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Convert amount to cents for Stripe
	amountInCents := int64(req.TotalAmount * 100)

	// Create Stripe PaymentIntent with only card payment method
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amountInCents),
		Currency:           stripe.String("usd"),
		PaymentMethodTypes: []*string{stripe.String("card")},
	}

	// Add metadata using SetMetadata
	params.AddMetadata("storeId", req.StoreId)
	params.AddMetadata("quantity", fmt.Sprintf("%d", req.Quantity))
	params.AddMetadata("pickup_time", req.PickupTime)

	pi, err := paymentintent.New(params)
	if err != nil {
		fmt.Println("Failed to create payment intent", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"clientSecret":    pi.ClientSecret,
		"paymentIntentId": pi.ID,
	})
}

func ConfirmReservation(c *gin.Context) {
	var req struct {
		PaymentIntentId string `json:"paymentIntentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Invalid request parameters", err)
		c.JSON(400, gin.H{"error": "Invalid request parameters"})
		return
	}

	// Verify payment status
	pi, err := paymentintent.Get(req.PaymentIntentId, nil)
	if err != nil {
		fmt.Println("Failed to verify payment", err)
		c.JSON(500, gin.H{"error": "Failed to verify payment"})
		return
	}

	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		fmt.Println("Payment not completed", pi.Status)
		c.JSON(400, gin.H{"error": "Payment not completed"})
		return
	}

	// Create reservation record in database
	_, err = db.DB.Exec(`
		INSERT INTO reservations (
			user_id, 
			store_id, 
			quantity, 
			total_amount, 
			status, 
			payment_id,
			pickup_time
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		c.GetString("user_id"),
		pi.Metadata["storeId"],
		parseInt(pi.Metadata["quantity"]),
		float64(pi.Amount)/100,
		"confirmed",
		pi.ID,
		pi.Metadata["pickup_time"],
	)

	if err != nil {
		fmt.Printf("Failed to create reservation record: %v\n", err)
		c.JSON(500, gin.H{"error": "Failed to create reservation record"})
		return
	}

	// Update bags_available count in stores table
	_, err = db.DB.Exec(`
		UPDATE stores 
		SET bags_available = GREATEST(0, bags_available - $1), updated_at = NOW()
		WHERE id = $2
	`, parseInt(pi.Metadata["quantity"]), pi.Metadata["storeId"])

	if err != nil {
		fmt.Printf("WARNING: Failed to update bags_available for store %s: %v\n", pi.Metadata["storeId"], err)
		// Don't fail the payment confirmation if bags_available update fails
	} else {
		fmt.Printf("Updated bags_available for store %s: decreased by %d\n", pi.Metadata["storeId"], parseInt(pi.Metadata["quantity"]))
	}

	// Create reservation record
	reservation := struct {
		StoreID     string  `json:"storeId"`
		UserID      string  `json:"userId"`
		Quantity    int     `json:"quantity"`
		TotalAmount float64 `json:"totalAmount"`
		Status      string  `json:"status"`
		PaymentID   string  `json:"paymentId"`
	}{
		StoreID:     pi.Metadata["storeId"],
		UserID:      c.GetString("userId"),
		Quantity:    parseInt(pi.Metadata["quantity"]),
		TotalAmount: float64(pi.Amount) / 100,
		Status:      "confirmed",
		PaymentID:   pi.ID,
	}

	c.JSON(200, gin.H{
		"status":      "success",
		"reservation": reservation,
	})
}

func ConfirmPayAtStore(c *gin.Context) {
	var req PayAtStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the payment intent to retrieve metadata
	pi, err := paymentintent.Get(req.PaymentIntentId, nil)
	if err != nil {
		fmt.Printf("Failed to retrieve payment intent: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment intent"})
		return
	}

	userID := c.GetString("user_id")
	storeID := pi.Metadata["storeId"]
	// fmt.Println("storeID", storeID)
	quantity := parseInt(pi.Metadata["quantity"])
	pickupTime := pi.Metadata["pickup_time"]
	// Get store price for total amount calculation
	var storePrice float64
	err = db.DB.Get(&storePrice, "SELECT price FROM stores WHERE id = $1", storeID)
	if err != nil {
		fmt.Println("Failed to get store details", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get store details"})
		return
	}

	totalAmount := storePrice * float64(quantity)

	// Insert the reservation
	_, err = db.DB.Exec(`
		INSERT INTO reservations 
		(user_id, store_id, quantity, total_amount, status, payment_id, pickup_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, storeID, quantity, totalAmount, "pending", "pay_at_store_"+req.PaymentIntentId, pickupTime)

	if err != nil {
		fmt.Printf("Failed to create reservation: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	// Update bags_available count in stores table
	_, err = db.DB.Exec(`
		UPDATE stores 
		SET bags_available = GREATEST(0, bags_available - $1), updated_at = NOW()
		WHERE id = $2
	`, quantity, storeID)

	if err != nil {
		fmt.Printf("WARNING: Failed to update bags_available for store %s: %v\n", storeID, err)
		// Don't fail the reservation creation if bags_available update fails
	} else {
		fmt.Printf("Updated bags_available for store %s: decreased by %d\n", storeID, quantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Reservation created successfully",
	})
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
