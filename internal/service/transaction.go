package service

import (
	"fmt"
	"github.com/google/uuid"
	database "go-ledger/internal/db"
	"go-ledger/pkg/models"
	"time"
	// Assuming you have a package that initializes and exports your DB connection
)

func CreateTransaction(req models.TransactionRequest) (*models.Transaction, error) {
	// Generate a new UUID for this transaction
	id := uuid.New().String()

	// Use the current timestamp for CreatedAt
	createdAt := time.Now()

	// Prepare the insert statement
	stmt := `INSERT INTO transaction (id, created_at, amount, action) VALUES ($1, $2, $3, $4)`

	// Execute the insert operation
	_, err := database.DB.Exec(stmt, id, createdAt, req.Amount, req.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to insert transaction: %w", err)
	}

	// Return the newly created transaction
	return &models.Transaction{
		ID:        id,
		CreatedAt: createdAt,
		Amount:    req.Amount,
		Action:    req.Action,
	}, nil
}
