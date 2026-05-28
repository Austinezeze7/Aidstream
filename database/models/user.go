package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID            string    `json:"id"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	Role          string    `json:"role"`
	WalletAddress string    `json:"wallet_address"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetAllUsers - fetch all users from database
func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, full_name, email, role, wallet_address, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.WalletAddress, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUserByEmail - find a user by email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id, full_name, email, password_hash, role, wallet_address, created_at FROM users WHERE email = $1", email).
		Scan(&u.ID, &u.FullName, &u.Email, &u.PasswordHash, &u.Role, &u.WalletAddress, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser - insert a new user
func CreateUser(db *sql.DB, u User) error {
	_, err := db.Exec(
		"INSERT INTO users (full_name, email, password_hash, role, wallet_address) VALUES ($1, $2, $3, $4, $5)",
		u.FullName, u.Email, u.PasswordHash, u.Role, u.WalletAddress,
	)
	return err
}