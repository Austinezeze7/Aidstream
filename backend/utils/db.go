package utils

import (
    "context"
    "fmt"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() error {
    connStr := os.Getenv("DATABASE_URL") // just the variable NAME
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