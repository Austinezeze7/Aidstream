package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	DB = db
	fmt.Println("✅ Connected to Neon database successfully!")
}

func RunMigrations() {
	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal("Error reading schema.sql: ", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal("Error running migrations: ", err)
	}

	fmt.Println("✅ Database tables created successfully!")
}
