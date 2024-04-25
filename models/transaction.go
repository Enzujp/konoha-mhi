package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	TransactionID uuid.UUID  `json:"transaction_id"`
	SenderID      string  	 `json:"sender_id"`
	ReceiverID    string  	 `json:"receiver_id"`
	Amount        float64 	 `json:"amount"`
	Timestamp     time.Time  `json:"timestamp"`
}