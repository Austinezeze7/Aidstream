package utils

import (
	"context"
	"fmt"
	"os"
	"time"

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
func GetCaseByID(id string) (map[string]interface{}, error) {
    row := DB.QueryRow(context.Background(),
        `SELECT id, title, description, category, target_amount, amount_raised, status
         FROM cases WHERE id = $1`, id,
    )

    var caseID, title, description, status string
    var category *string
    var targetAmount, amountRaised float64

    err := row.Scan(&caseID, &title, &description, &category, &targetAmount, &amountRaised, &status)
    if err != nil {
        return nil, fmt.Errorf("case not found: %v", err)
    }

    c := map[string]interface{}{
        "id":            caseID,
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

    return c, nil
}
func CreateDonation(caseID, donorID, paymentMethod string, amount float64) error {
	tx, err := DB.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// donorID can be empty — store as NULL if so
	var donorParam interface{}
	if donorID == "" {
		donorParam = nil
	} else {
		donorParam = donorID
	}

	// Insert into donations — status defaults to 'pending'
	_, err = tx.Exec(context.Background(),
		`INSERT INTO donations (case_id, donor_id, amount, status)
		 VALUES ($1, $2, $3, 'pending')`,
		caseID, donorParam, amount,
	)
	if err != nil {
		return fmt.Errorf("donation insert failed: %v", err)
	}

	// Update amount_raised on the case
	_, err = tx.Exec(context.Background(),
		`UPDATE cases SET amount_raised = amount_raised + $1 WHERE id = $2`,
		amount, caseID,
	)
	if err != nil {
		return fmt.Errorf("case update failed: %v", err)
	}

	return tx.Commit(context.Background())
}
func LoginUser(email, password string) (map[string]interface{}, error) {
	row := DB.QueryRow(context.Background(),
		`SELECT id, email, role, first_name FROM users 
		 WHERE email = $1 AND password = $2`,
		email, password,
	)

	var id, userEmail, role, firstName string
	err := row.Scan(&id, &userEmail, &role, &firstName)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return map[string]interface{}{
		"id":         id,
		"email":      userEmail,
		"role":       role,
		"first_name": firstName,
	}, nil
}
// =========================
// ADMIN STATS
// =========================
func GetAdminStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	// Pending verifications (cases awaiting review)
	var pendingVerifications int
	DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM cases WHERE status = 'pending'`,
	).Scan(&pendingVerifications)

	// Pending disbursements
	var pendingDisbursements int
	DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM donations WHERE status = 'pending'`,
	).Scan(&pendingDisbursements)

	// Active cases
	var activeCases int
	DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM cases WHERE status = 'active'`,
	).Scan(&activeCases)

	// Completed cases
	var completedCases int
	DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM cases WHERE status = 'closed'`,
	).Scan(&completedCases)

	stats["pending_verifications"] = pendingVerifications
	stats["pending_disbursements"] = pendingDisbursements
	stats["active_cases"]          = activeCases
	stats["completed_cases"]       = completedCases

	return stats, nil
}

// =========================
// PENDING VERIFICATIONS
// =========================
func GetPendingVerifications() ([]map[string]interface{}, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT c.id, c.title, c.category, c.created_at,
		        u.first_name, u.last_name
		 FROM cases c
		 LEFT JOIN users u ON c.created_by = u.id
		 WHERE c.status = 'pending'
		 ORDER BY c.created_at DESC
		 LIMIT 10`,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, title string
		var category, firstName, lastName *string
		var createdAt time.Time  // ← time.Time not string

		err := rows.Scan(&id, &title, &category, &createdAt, &firstName, &lastName)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}

		name := "Unknown"
		if firstName != nil && lastName != nil {
			name = *firstName + " " + *lastName
		}

		cat := ""
		if category != nil {
			cat = *category
		}

		list = append(list, map[string]interface{}{
			"id":         id,
			"title":      title,
			"category":   cat,
			"created_at": createdAt.Format(time.RFC3339),  // ← format to string
			"name":       name,
		})
	}
	return list, nil
}

// =========================
// RECENT DISBURSEMENTS
// =========================
func GetRecentDisbursements() ([]map[string]interface{}, error) {
	rows, err := DB.Query(context.Background(),
		`SELECT d.id, d.amount, d.status, d.donated_at,
		        COALESCE(d.transaction_hash, 'N/A'),
		        c.title,
		        u.first_name, u.last_name
		 FROM donations d
		 LEFT JOIN cases c ON d.case_id = c.id
		 LEFT JOIN users u ON d.donor_id = u.id
		 ORDER BY d.donated_at DESC
		 LIMIT 10`,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, status, txHash, caseTitle string
		var amount float64
		var donatedAt time.Time  // ← time.Time not string
		var firstName, lastName *string

		err := rows.Scan(&id, &amount, &status, &donatedAt, &txHash, &caseTitle, &firstName, &lastName)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}

		name := "Anonymous"
		if firstName != nil && lastName != nil {
			name = *firstName + " " + *lastName
		}

		list = append(list, map[string]interface{}{
			"id":               id,
			"amount":           amount,
			"status":           status,
			"donated_at":       donatedAt.Format(time.RFC3339),  // ← format to string
			"transaction_hash": txHash,
			"case_title":       caseTitle,
			"donor_name":       name,
		})
	}
	return list, nil
}