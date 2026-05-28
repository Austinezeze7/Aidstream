package models

import (
	"database/sql"
	"time"
)

// User represents the updated users table schema
type User struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"first_name"`
	MiddleName    *string   `json:"middle_name"` // Pointer handles nullable database fields
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"` // Hides password field from JSON responses for security
	Role          string    `json:"role"`
	WalletAddress *string   `json:"wallet_address"`
	CreatedAt     time.Time `json:"created_at"`
}

// FindUserByEmail fetches a user for the login validation mechanism
func FindUserByEmail(db *sql.DB, email string) (*User, error) {
	var user User
	query := `
		SELECT id, first_name, middle_name, last_name, email, password, role, wallet_address, created_at 
		FROM users 
		WHERE email = $1`

	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.WalletAddress,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser inserts a registration profile with the new decoupled name structure
func CreateUser(db *sql.DB, firstName, middleName, lastName, email, password, role string) (*User, error) {
	var user User
	query := `
		INSERT INTO users (first_name, middle_name, last_name, email, password, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, first_name, middle_name, last_name, email, role, created_at`

	// Handle optional middle name values safely
	var dbMiddleName interface{} = nil
	if middleName != "" {
		dbMiddleName = middleName
	}

	err := db.QueryRow(query, firstName, dbMiddleName, lastName, email, password, role).Scan(
		&user.ID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers fetches all profiles recorded in the system
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `SELECT id, first_name, middle_name, last_name, email, role, wallet_address, created_at FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.MiddleName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.WalletAddress,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
