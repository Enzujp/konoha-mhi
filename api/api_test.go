package api_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"github.com/google/uuid"
    "github.com/stretchr/testify/assert"
	"fmt"
	"github.com/go-chi/chi/v5"
	
)

type User struct {
    ID            string  `json:"ID"`
    FirstName     string  `json:"FirstName"`
    LastName      string  `json:"LastName"`
    Email         string  `json:"Email"`
    WalletBalance float64 `json:"WalletBalance"`
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
    // Extract userID from the request URL
    userID := r.URL.Path[len("/users/"):]

    // Simulate database behavior
    var user *User
    switch userID {
    case "1":
        user = &User{
            ID:            "1",
            FirstName:     "Toph",
            LastName:      "Beifong",
            Email:         "Toph.b@example.com",
            WalletBalance: 100.00,
        }
    default:
        // Return 404 Not Found if user is not found
        w.WriteHeader(http.StatusNotFound)
        return
    }

    // Marshal user details to JSON and write response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func TestGetUserDetailsExistingUser(t *testing.T) {
    req, err := http.NewRequest("GET", "/users/1", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    GetUserDetails(rr, req)

    // Check if the response status code is OK
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check if the response body contains the expected user details
    expectedBody := `{"ID":"1","FirstName":"Toph","LastName":"Beifong","Email":"Toph.b@example.com","WalletBalance":100}`
    responseBody := rr.Body.String()
    if responseBody != expectedBody {
        t.Errorf("handler returned unexpected body: got %v want %v", responseBody, expectedBody)
    }

    // Check if the response Content-Type header is set to application/json
    contentType := rr.Header().Get("Content-Type")
    if contentType != "application/json" {
        t.Errorf("incorrect Content-Type header: got %v want application/json", contentType)
    }
}

func TestGetUserDetailsNonExistingUser(t *testing.T) {
    req, err := http.NewRequest("GET", "/users/999", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    GetUserDetails(rr, req)

    // Check if the response status code is Not Found
    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("incorrect status code response: got %v want %v", status, http.StatusNotFound)
    }

    // Check if the response body contains the expected error message
    expectedBody := "User not found"
    responseBody := rr.Body.String()
    if responseBody != expectedBody {
        t.Errorf("iuncorrect body response : got %v want %v", responseBody, expectedBody)
    }
}



// Test TransactionDetails
type Transaction struct {
    TransactionID string `json:"transactionID"`
    SenderID      string `json:"senderID"`
    ReceiverID    string `json:"receiverID"`
    Amount        int    `json:"amount"`
    Timestamp     string `json:"timestamp"`
}

func GetTransactionDetails(w http.ResponseWriter, r *http.Request) {
    var transaction *Transaction
    idString := chi.URLParam(r, "transactionID")
    _, err := uuid.Parse(idString)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, `{"error": "invalid ID"}`)
        return
    }
    // Simulate database behavior
    switch idString {
    case "1":
        transaction = &Transaction{
            TransactionID: "1",
            SenderID:      "sender_id_1",
            ReceiverID:    "receiver_id_1",
            Amount:        100,
            Timestamp:     "2022-04-27T10:00:00Z",
        }
    case "2":
        transaction = &Transaction{
            TransactionID: "2",
            SenderID:      "sender_id_2",
            ReceiverID:    "receiver_id_2",
            Amount:        200,
            Timestamp:     "2022-04-27T11:00:00Z",
        }
    default:
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, `{"error": "transaction not found"}`)
        return
    }

    // Write response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transaction)
}

func TestGetTransactionDetailsExistingTransaction(t *testing.T) {
    req, err := http.NewRequest("GET", "/transactions/1", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    GetTransactionDetails(rr, req)

    // Check if the response status code is OK
    assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

    // Check if the response body contains the expected transaction details
    expectedBody := `{"transactionID":"1","senderID":"sender_id_1","receiverID":"receiver_id_1","amount":100,"timestamp":"2022-04-27T10:00:00Z"}`
    assert.Equal(t, expectedBody, rr.Body.String(), "handler returned unexpected body")
}

func TestGetTransactionDetailsNonExistingTransaction(t *testing.T) {
    req, err := http.NewRequest("GET", "/transactions/999", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    GetTransactionDetails(rr, req)

    // Check if the response status code is Not Found
    assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")

    // Check if the response body contains the expected error message
    expectedBody := `{"error": "transaction not found"}`
    assert.Equal(t, expectedBody, rr.Body.String(), "handler returned unexpected body")
}

func TestGetTransactionDetailsInvalidID(t *testing.T) {
    req, err := http.NewRequest("GET", "/transactions/invalid", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    GetTransactionDetails(rr, req)

    // Check if the response status code is Bad Request
    assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")

    // Check if the response body contains the expected error message
    expectedBody := `{"error": "invalid ID"}`
    assert.Equal(t, expectedBody, rr.Body.String(), "handler returned unexpected body")
}