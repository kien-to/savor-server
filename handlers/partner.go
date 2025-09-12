package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"savor-server/db"

	"github.com/gin-gonic/gin"
)

// PartnerContactRequest represents the partner contact form data
type PartnerContactRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	StoreName string `json:"storeName" binding:"required"`
	Message   string `json:"message"`
}

// PartnerContact represents the database model for partner contacts
type PartnerContact struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	StoreName string    `json:"store_name" db:"store_name"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Status    string    `json:"status" db:"status"` // 'new', 'contacted', 'completed'
}

// @Summary Submit partner contact form
// @Description Submit a contact form for potential store partners
// @Tags partner
// @Accept json
// @Produce json
// @Param request body PartnerContactRequest true "Partner contact information"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/partner/contact [post]
func SubmitPartnerContact(c *gin.Context) {
	var req PartnerContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding partner contact request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create the partner_contacts table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS partner_contacts (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			phone VARCHAR(50) NOT NULL,
			store_name VARCHAR(255) NOT NULL,
			message TEXT,
			status VARCHAR(20) DEFAULT 'new',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.DB.Exec(createTableQuery)
	if err != nil {
		log.Printf("Error creating partner_contacts table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database initialization error",
		})
		return
	}

	// Insert the contact form data
	insertQuery := `
		INSERT INTO partner_contacts (name, email, phone, store_name, message, status, created_at)
		VALUES ($1, $2, $3, $4, $5, 'new', CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`

	var contactID int
	var createdAt time.Time
	err = db.DB.QueryRow(insertQuery, req.Name, req.Email, req.Phone, req.StoreName, req.Message).Scan(&contactID, &createdAt)
	if err != nil {
		log.Printf("Error inserting partner contact: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save contact information",
		})
		return
	}

	log.Printf("New partner contact submitted: ID=%d, Name=%s, Email=%s, Store=%s",
		contactID, req.Name, req.Email, req.StoreName)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Contact form submitted successfully",
		"contact_id": contactID,
		"created_at": createdAt,
	})
}

// @Summary Get all partner contacts (Admin only)
// @Description Retrieve all partner contact submissions for admin review
// @Tags partner
// @Produce json
// @Success 200 {array} PartnerContact "List of partner contacts"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/partner/contacts [get]
func GetPartnerContacts(c *gin.Context) {
	query := `
		SELECT id, name, email, phone, store_name, message, status, created_at
		FROM partner_contacts
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying partner contacts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve contacts",
		})
		return
	}
	defer rows.Close()

	var contacts []PartnerContact
	for rows.Next() {
		var contact PartnerContact
		err := rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.Email,
			&contact.Phone,
			&contact.StoreName,
			&contact.Message,
			&contact.Status,
			&contact.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning partner contact: %v", err)
			continue
		}
		contacts = append(contacts, contact)
	}

	c.JSON(http.StatusOK, gin.H{
		"contacts": contacts,
		"total":    len(contacts),
	})
}

// @Summary Update partner contact status
// @Description Update the status of a partner contact (Admin only)
// @Tags partner
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Param request body map[string]string true "Status update"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/partner/contacts/{id}/status [put]
func UpdatePartnerContactStatus(c *gin.Context) {
	contactID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"new":       true,
		"contacted": true,
		"completed": true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status. Must be one of: new, contacted, completed",
		})
		return
	}

	updateQuery := `
		UPDATE partner_contacts 
		SET status = $1 
		WHERE id = $2
		RETURNING id, status
	`

	var updatedID int
	var updatedStatus string
	err := db.DB.QueryRow(updateQuery, req.Status, contactID).Scan(&updatedID, &updatedStatus)
	if err != nil {
		log.Printf("Error updating partner contact status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update contact status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Contact status updated to %s", updatedStatus),
		"id":      updatedID,
		"status":  updatedStatus,
	})
}
