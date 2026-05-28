package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// =========================
// CONNECT DATABASE
// =========================
func ConnectDB() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %v", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully")
	return nil
}

// =========================
// CREATE USER
// =========================
func CreateUser(firstName, middleName, lastName, email, password, role, walletAddress string) error {
	_, err := DB.Exec(context.Background(),
		`INSERT INTO users (
			email,
			role,
			wallet_address,
			first_name,
			middle_name,
			last_name,
			password
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		email,
		role,
		walletAddress,
		firstName,
		middleName,
		lastName,
		password,
	)
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}
	return nil
}

// =========================
// GET CASES
// =========================
func GetCases() ([]map[string]interface{}, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT id, title, description, category, target_amount, amount_raised, status
		 FROM cases
		 WHERE status = 'active'`,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var cases []map[string]interface{}

	for rows.Next() {
		var id, status, title, description string
		var category *string  // nullable
		var targetAmount, amountRaised float64

		err := rows.Scan(&id, &title, &description, &category, &targetAmount, &amountRaised, &status)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}

		c := map[string]interface{}{
			"id":            id,
			"title":         title,
			"description":   description,
			"category":      "",
			"target_amount": targetAmount,
			"amount_raised": amountRaised,
			"status":        status,
		}

		if category != nil {
			c["category"] = *category
		}

		cases = append(cases, c)
	}

	return cases, nil
}