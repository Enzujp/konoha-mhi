package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"encoding/json"

	"github.com/enzujp/konoha-mhi/database"
	"github.com/enzujp/konoha-mhi/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (a *API) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	// Obtain userID from request parameters
	id := chi.URLParam(r, "userID")

	var user models.User
	// Search database for user
	err := database.DB.QueryRow("SELECT id, first_name, last_name, email, wallet_balance FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.WalletBalance)
	if err != nil {
		// If row does not exist
		if err == sql.ErrNoRows {
			renderStatusError(r, w, http.StatusNotFound, "User not found")
			return
		}
		renderStatusError(r, w, http.StatusInternalServerError, "Internal server error")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)	
}


type Disbursement struct {
	ReceiverID string  `json:"user_id"`
	Amount     float64 `json:"amount"`
}

func (a *API) DisburseFunds(w http.ResponseWriter, r *http.Request) {
	var requestBody []Disbursement

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		renderStatusError(r, w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Retrieve sender ID from URL parameter
	senderID, err := getSenderIDFromRequest(r)
	if err != nil {
		renderStatusError(r, w, http.StatusBadRequest, ErrInvalidUserID)
		return
	}

	// Check sender's balance
	senderBalance, err := getSenderBalance(senderID)
	if err != nil {
		renderStatusError(r, w, http.StatusInternalServerError, ErrFailedToFetchBalance)
		return
	}

	totalAmount := 0.0

	// Validate disbursements and calculate total amount
	for _, d := range requestBody {
		if d.Amount <= 0 {
			renderStatusError(r, w, http.StatusBadRequest, ErrInvalidAmount)
			return
		}
		totalAmount += d.Amount
	}

	if senderBalance < totalAmount {
		renderStatusError(r, w, http.StatusBadRequest, ErrInsufficientBalance)
		return
	}

	// Start a database transaction
	tx, err := database.DB.Begin()
	if err != nil {
		renderStatusError(r, w, http.StatusInternalServerError, ErrFailedToStartTransaction)
		return
	}
	defer tx.Rollback()

	// Process disbursements and log transactions
	transactionIDs := make(map[string]string)
	for _, d := range requestBody {
		transactionID := uuid.New().String()
		// Update receiver's balance
		_, err := tx.Exec("UPDATE users SET wallet_balance = wallet_balance + $1 WHERE id = $2", d.Amount, d.ReceiverID)
		if err != nil {
			renderStatusError(r, w, http.StatusInternalServerError, ErrFailedToUpdateReceiver)
			return
		}

		// Insert transaction record into database
		_, err = tx.Exec("INSERT INTO transactions (id, sender_id, receiver_id, amount, timestamp) VALUES($1, $2, $3, $4, $5)", transactionID, senderID, d.ReceiverID, d.Amount, time.Now())
		if err != nil {
			renderStatusError(r, w, http.StatusInternalServerError, ErrFailedToLogTransaction)
			return
		}
		// Populate map with transactionID
		transactionIDs[d.ReceiverID] = transactionID
	}

	// Update sender's balance
	_, err = tx.Exec("UPDATE users SET wallet_balance = wallet_balance - $1 WHERE id = $2", totalAmount, senderID)
	if err != nil {
		renderStatusError(r, w, http.StatusInternalServerError, ErrFailedToUpdateSender)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		renderStatusError(r, w, http.StatusInternalServerError, ErrFailedToCommit)
		return
	}

	// Respond with success message and total amount disbursed
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"message":              "Funds disbursed successfully!",
		"total_amount":         totalAmount,
		"total_disbursements":  len(requestBody),
		"disbursements":        requestBody,
		"balance_before":       senderBalance,
		"balance_after":        senderBalance - totalAmount,
		"transaction_id":       transactionIDs,
	})
}

func getSenderBalance(senderID uuid.UUID) (float64, error) {
	var senderBalance float64
	err := database.DB.QueryRow("SELECT wallet_balance FROM users WHERE id = $1", senderID).Scan(&senderBalance)
	return senderBalance, err
}

func getSenderIDFromRequest(r *http.Request) (uuid.UUID, error) {
	idString := chi.URLParam(r, "userID")
	return uuid.Parse(idString)
}

// Fetches transaction details using transactionID
func (a *API) GetTransactionDetails(w http.ResponseWriter, r *http.Request) {
    var transaction models.Transaction
    idString := chi.URLParam(r, "transactionID")
    id, err := uuid.Parse(idString)
    if err != nil{
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, map[string]interface{}{"error": "invalid ID"})
        return
    }
    err = database.DB.QueryRow("SELECT id, sender_id, receiver_id, amount, timestamp FROM transactions WHERE id = $1", id).Scan(&transaction.TransactionID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Timestamp)
    if err != nil {
        fmt.Println("This is the error, here?", err)
        render.Status(r, http.StatusInternalServerError)
        render.JSON(w, r, map[string]interface{}{"error": "internal server error"})
        return
    }
    render.Status(r, http.StatusOK)
    render.JSON(w, r, map[string]interface{}{
        "transaction_details": transaction,
    })
}

