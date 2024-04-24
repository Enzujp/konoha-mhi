package models

type User struct {
	ID	string	`json:"id"`
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
	Email	string	`json:"email"`
	WalletBalance	float64	`json:"wallet_balance"`
}