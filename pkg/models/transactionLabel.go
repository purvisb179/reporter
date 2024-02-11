package models

type TransactionLabel struct {
	TransactionID int64 `json:"transaction_id" db:"transaction_id"`
	LabelID       int64 `json:"label_id" db:"label_id"`
}
