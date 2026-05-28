package models

import (
	"database/sql"
	"time"
)

type Case struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	TargetAmount float64   `json:"target_amount"`
	AmountRaised float64   `json:"amount_raised"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetAllCases - fetch all cases
func GetAllCases(db *sql.DB) ([]Case, error) {
	rows, err := db.Query("SELECT id, title, description, category, target_amount, amount_raised, status, created_by, created_at FROM cases")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cases []Case
	for rows.Next() {
		var c Case
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Category, &c.TargetAmount, &c.AmountRaised, &c.Status, &c.CreatedBy, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, nil
}

// GetCaseByID - fetch a single case
func GetCaseByID(db *sql.DB, id string) (*Case, error) {
	var c Case
	err := db.QueryRow("SELECT id, title, description, category, target_amount, amount_raised, status, created_by, created_at FROM cases WHERE id = $1", id).
		Scan(&c.ID, &c.Title, &c.Description, &c.Category, &c.TargetAmount, &c.AmountRaised, &c.Status, &c.CreatedBy, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// CreateCase - insert a new case
func CreateCase(db *sql.DB, c Case) error {
	_, err := db.Exec(
		"INSERT INTO cases (title, description, category, target_amount, created_by) VALUES ($1, $2, $3, $4, $5)",
		c.Title, c.Description, c.Category, c.TargetAmount, c.CreatedBy,
	)
	return err
}

// UpdateCaseStatus - change case status
func UpdateCaseStatus(db *sql.DB, id string, status string) error {
	_, err := db.Exec("UPDATE cases SET status = $1 WHERE id = $2", status, id)
	return err
}