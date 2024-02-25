package service

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type ReportService struct {
	DB *sql.DB
}

func NewReportService(db *sql.DB) *ReportService {
	return &ReportService{DB: db}
}

// GenerateReport generates an Excel report for transactions and their labels.
func (rs *ReportService) GenerateReport(labels []string) (*xlsx.File, error) {
	// Start with the base SQL query
	query := `
SELECT t.amount, l.value, t.created_at
FROM transaction t
JOIN transaction_label tl ON t.id = tl.transaction_id
JOIN label l ON tl.label_id = l.id
`

	// Modify the query to filter by labels if they are provided
	var params []interface{}
	if len(labels) > 0 {
		// PostgreSQL uses numbered placeholders, e.g., $1, $2, etc.
		placeholder := make([]string, len(labels))
		for i := range labels {
			placeholder[i] = fmt.Sprintf("$%d", i+1)
			params = append(params, labels[i])
		}
		query += " WHERE l.value IN (" + strings.Join(placeholder, ",") + ")"
	}

	// Prepare the query with parameters to avoid SQL injection
	stmt, err := rs.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}
	defer stmt.Close()

	// Execute the query with parameters
	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	// Create a new Excel file and a sheet
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Report")
	if err != nil {
		return nil, fmt.Errorf("failed to add sheet: %w", err)
	}

	// Add headers to the Excel sheet
	headers := []string{"Amount", "Label Value", "Created At"}
	headerRow := sheet.AddRow()
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	// Iterate through the query results
	for rows.Next() {
		var amount int
		var value string
		var createdAt time.Time

		err := rows.Scan(&amount, &value, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Add rows to the sheet for each record
		row := sheet.AddRow()
		row.AddCell().SetInt(amount)
		row.AddCell().SetValue(value)
		row.AddCell().SetDateTime(createdAt)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Return the generated Excel file
	return file, nil
}

// GetDistinctLabelValues retrieves a distinct list of label values from the labels table.
func (rs *ReportService) GetDistinctLabelValues() ([]string, error) {
	// Define the SQL query to select distinct label values
	query := `SELECT DISTINCT value FROM label ORDER BY value;`

	// Execute the query
	rows, err := rs.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	// Slice to hold the distinct label values
	var values []string

	// Iterate through the query results
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		values = append(values, value)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Return the distinct label values
	return values, nil
}
