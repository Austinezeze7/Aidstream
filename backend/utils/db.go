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