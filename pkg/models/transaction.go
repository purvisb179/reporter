package models

import "time"

type Transaction struct {
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Amount    int       `json:"amount" db:"amount"`
}

type TransactionRequest struct {
	Amount int `json:"amount" binding:"required"`
}

// labels

type Label struct {
	ID    int64  `json:"id" db:"id"`
	Key   string `json:"key" db:"key"`
	Value string `json:"value" db:"value"`
}

// transaction label

type TransactionLabel struct {
	TransactionID int64 `json:"transaction_id" db:"transaction_id"`
	LabelID       int64 `json:"label_id" db:"label_id"`
}
