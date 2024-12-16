package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type ReservationRequest struct {
	StoreId       string  `json:"storeId" binding:"required"`
	Quantity      int     `json:"quantity" binding:"required,min=1"`
	TotalAmount   float64 `json:"totalAmount" binding:"required"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
}

func CreateReservation(c *gin.Context) {
	var req ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request parameters"})
		return
	}

	// Convert amount to cents for Stripe
	amountInCents := int64(req.TotalAmount * 100)

	// Create Stripe PaymentIntent with only card payment method
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amountInCents),
		Currency:          stripe.String("usd"),
		PaymentMethodTypes: []*string{stripe.String("card")},
	}

	// Add metadata using SetMetadata
	params.AddMetadata("storeId", req.StoreId)
	params.AddMetadata("quantity", fmt.Sprintf("%d", req.Quantity))

	pi, err := paymentintent.New(params)
	if err != nil {
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
		c.JSON(400, gin.H{"error": "Invalid request parameters"})
		return
	}

	// Verify payment status
	pi, err := paymentintent.Get(req.PaymentIntentId, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to verify payment"})
		return
	}

	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		c.JSON(400, gin.H{"error": "Payment not completed"})
		return
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

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
