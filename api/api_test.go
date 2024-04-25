package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enzujp/konoha-mhi/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserDetails_UserFound(t *testing.T) {
	// Mocking a user ID for the test
	userID := uuid.New()

	// Mocking a user based on the provided JSON
	expectedUser := models.User{
		ID:            userID,
		FirstName:     "Sukuna",
		LastName:      "Shinazugawa",
		Email:         "sukuna.shinazugawa@example.com",
		WalletBalance: 18770,
	}

	// Mocking the request
	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)

	// Mocking the response recorder
	rr := httptest.NewRecorder()

	// Creating an instance of your API handler
	a := &API{}

	// Calling the handler function
	a.GetUserDetails(rr, req)

	// Checking the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parsing the response body into a User struct
	var userResponse models.User
	err := json.NewDecoder(rr.Body).Decode(&userResponse)
	if err != nil {
		t.Errorf("Error decoding JSON response: %v", err)
	}

	// Checking if the retrieved user matches the expected user
	assert.Equal(t, expectedUser, userResponse)
}

// func TestGetUserDetails_UserNotFound(t *testing.T) {
// 	// Mocking a non-existent user ID for the test
// 	userID := "nonexistent-user-id"

// 	// Mocking the request
// 	req := httptest.NewRequest("GET", "/users/"+userID, nil)

// 	// Mocking the response recorder
// 	rr := httptest.NewRecorder()

// 	// Creating an instance of your API handler
// 	a := &API{}

// 	// Calling the handler function
// 	a.GetUserDetails(rr, req)

// 	// Checking the response status code
// 	assert.Equal(t, http.StatusNotFound, rr.Code)
// }
