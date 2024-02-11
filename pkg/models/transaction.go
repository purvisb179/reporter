package models

import "time"

type ActionType string

const (
	ActionAdd    ActionType = "add"
	ActionSub    ActionType = "sub"
	ActionAdjust ActionType = "adj"
)

type Transaction struct {
	ID        string     `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Amount    int        `json:"amount" db:"amount"`
	Action    ActionType `json:"action" db:"action"`
}

// TransactionRequest represents the request payload for creating a transaction.
type TransactionRequest struct {
	Amount int        `json:"amount" binding:"required"`
	Action ActionType `json:"action" binding:"required,oneof=add sub adj"`
}
