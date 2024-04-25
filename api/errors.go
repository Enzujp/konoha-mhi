package api

const (
	ErrInvalidUserID            = "invalid user ID"
	ErrInvalidRequestBody       = "invalid request body"
	ErrFailedToFetchBalance     = "failed to fetch sender's balance"
	ErrInvalidAmount            = "invalid amount for disbursement"
	ErrInsufficientBalance      = "insufficient balance to carry out transaction"
	ErrFailedToStartTransaction = "failed to start transaction"
	ErrFailedToUpdateReceiver   = "failed to update receiver's balance"
	ErrFailedToLogTransaction   = "failed to log transaction"
	ErrFailedToUpdateSender     = "failed to update sender's balance"
	ErrFailedToCommit           = "failed to commit transaction"
)