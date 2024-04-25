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

type BulkDisbursementRequest struct {
    Receivers []models.User `json:"receivers"`
    Amount    float64       `json:"amount"`
}

func (a *API) DisburseFunds(w http.ResponseWriter, r *http.Request) {
    var requestBody BulkDisbursementRequest
	// Obtain sender's ID
    idString := chi.URLParam(r, "id")
    senderID, err := uuid.Parse(idString)
    if err != nil {
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, map[string]interface{}{"error": "invalid user ID"})
        return
    }

    // Check sender's balance
    var senderBalance float64
    if err := database.DB.QueryRow("SELECT wallet_balance FROM users WHERE id = $1", senderID).Scan(&senderBalance); err != nil {
        render.Status(r, http.StatusInternalServerError)
        render.JSON(w, r, map[string]interface{}{"error": "internal server error"})
        return
    }

    if senderBalance < requestBody.Amount*float64(len(requestBody.Receivers)) {
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, map[string]interface{}{"error": "insufficient balance to carry out transaction"})
        return
    }

    // Start a database transaction
    tx, err := database.DB.Begin()
    if err != nil {
        render.Status(r, http.StatusInternalServerError)
        render.JSON(w, r, map[string]interface{}{"error": "internal server error"})
        return
    }
    defer tx.Rollback()

    // Process receivers
    for _, receiver := range requestBody.Receivers {
        // Update receiver's balance
        receiver.WalletBalance += requestBody.Amount
        _, err := tx.Exec("UPDATE users SET wallet_balance = $1 WHERE id = $2", receiver.WalletBalance, receiver.ID)
        if err != nil {
            render.Status(r, http.StatusInternalServerError)
            render.JSON(w, r, map[string]interface{}{"error": "internal server error"})
            return
        }
    }

    // Update sender's balance
    senderBalance -= requestBody.Amount * float64(len(requestBody.Receivers))
    _, err = tx.Exec("UPDATE users SET wallet_balance = $1 WHERE id = $2", senderBalance, senderID)
    if err != nil {
        render.Status(r, http.StatusInternalServerError)
        render.JSON(w, r, map[string]interface{}{"error": "internal server error"})
        return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        render.Status(r, http.StatusInternalServerError)
        render.JSON(w, r, map[string]interface{}{"error": "internal server error"})
        return
    }

    // Respond with success message
    render.Status(r, http.StatusOK)
    render.JSON(w, r, map[string]interface{}{
		"message": "Funds disbursed successfully!",
		"total_amount":	totalAmount,
		"total_disbursements":	len(requestBody.Disbursements),
	})
}