package models

import "time"

type ActionType string

const (
	ActionAdd    ActionType = "add"
	ActionSub    ActionType = "sub"
	ActionAdjust ActionType = "adj"
)

type Transaction struct {
	ID        int64      `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Amount    float64    `json:"amount" db:"amount"`
	Action    ActionType `json:"action" db:"action"`
}
