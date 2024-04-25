package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/enzujp/konoha-mhi/database"
	"github.com/enzujp/konoha-mhi/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (a *API) GetUserDetails(w http.ResponseWriter, r *http.Request){
	// Obtain userID from request parameters
	idString := chi.URLParam(r, "id")
	id, err := uuid.Parse(idString)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]interface{}{"error": "Invalid ID"})
		return
	}

	var user models.User
	// Search database for user
	err = database.DB.QueryRow("SELECT id, first_name, last_name, email, wallet_balance FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.WalletBalance)
	if err != nil {
		// If row does not exist
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{"error": "User not found"})
			return
		}
		fmt.Println("Error querying row:", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, (map[string]interface{}{"error": "internal server error"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type DisbursementRequest struct {
	SenderID     string        `json:"sender_id"`
	Disbursements []Disbursement `json:"disbursements"`
}

type Disbursement struct {
	ReceiverID string  `json:"receiver_id"`
	Amount     float64 `json:"amount"`
}

func (a *API) DisburseFunds(w http.ResponseWriter, r *http.Request) {
	var requestBody DisbursementRequest
	idString := chi.URLParam(r, "id") // Retrieve sender ID from URL parameter
	senderID, err := uuid.Parse(idString)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]interface{}{"error": "invalid user ID"})
		return
	}

	// Check sender's balance
	var senderBalance float64
	err = database.DB.QueryRow("SELECT wallet_balance FROM users WHERE id = $1", senderID).Scan(&senderBalance)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{"error": "failed to fetch sender's balance"})
		return
	}

	totalAmount := 0.0

	// Validate disbursements and calculate total amount
	for _, d := range requestBody.Disbursements {
		if d.Amount <= 0 {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]interface{}{"error": "invalid amount for disbursement"})
			return
		}
		totalAmount += d.Amount
	}

	if senderBalance < totalAmount {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]interface{}{"error": "insufficient balance to carry out transaction"})
		return
	}

	// Start a database transaction
	tx, err := database.DB.Begin()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{"error": "failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Process disbursements
	for _, d := range requestBody.Disbursements {
		// Update receiver's balance
		_, err := tx.Exec("UPDATE users SET wallet_balance = wallet_balance + $1 WHERE id = $2", d.Amount, d.ReceiverID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{"error": "failed to update receiver's balance"})
			return
		}
	}

	// Update sender's balance
	_, err = tx.Exec("UPDATE users SET wallet_balance = wallet_balance - $1 WHERE id = $2", totalAmount, senderID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{"error": "failed to update sender's balance"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{"error": "failed to commit transaction"})
		return
	}

	// Respond with success message
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
        "message": "Funds disbursed successfully!",
        "total_amount": totalAmount,
        "total_disbursements": len(requestBody.Disbursements),
    })
}
