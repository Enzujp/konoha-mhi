package api

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	fmt.Println(err)
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
		render.JSON(w, r, (map[string]interface{}{"error": "Internal Server Error"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}


func (a *API) DisburseFunds(w http.ResponseWriter, r *http.Request){
	var users []models.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if err := disburse(user.ID, user.WalletBalance); err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Funds disbursed successfully"))
	// perhaps show user balance
}

func disburse(userID string, amount float64) error {
	// Get User's wallet balance
	var currentBalance float64
	err := database.DB.QueryRow("SELECT wallet_balance FROM users WHERE id = $1", userID).Scan(&currentBalance)
	if err != nil {
		return err
	}
	
	// Check to confirm that user has sufficient balance to carry out transaction
	if currentBalance < amount {
		return errors.New("insufficient balance, cannot complete this transaction")
	}
	// Update wallet balance
	_, err = database.DB.Exec("UPDATE users SET wallet_balance = wallet_balance - $1 WHERE id = $2", amount, userID)
	if err != nil {
		return err
	}

	return nil
}