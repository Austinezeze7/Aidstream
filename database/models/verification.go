package models

import (
	"database/sql"
	"time"
)

type Verification struct {
	ID         string    `json:"id"`
	CaseID     string    `json:"case_id"`
	VerifierID string    `json:"verifier_id"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
	VerifiedAt time.Time `json:"verified_at"`
}

// GetVerificationsByCase - fetch all verifications for a case
func GetVerificationsByCase(db *sql.DB, caseID string) ([]Verification, error) {
	rows, err := db.Query("SELECT id, case_id, verifier_id, status, notes, verified_at FROM verifications WHERE case_id = $1", caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verifications []Verification
	for rows.Next() {
		var v Verification
		err := rows.Scan(&v.ID, &v.CaseID, &v.VerifierID, &v.Status, &v.Notes, &v.VerifiedAt)
		if err != nil {
			return nil, err
		}
		verifications = append(verifications, v)
	}
	return verifications, nil
}

// CreateVerification - insert a new verification
func CreateVerification(db *sql.DB, v Verification) error {
	_, err := db.Exec(
		"INSERT INTO verifications (case_id, verifier_id, notes) VALUES ($1, $2, $3)",
		v.CaseID, v.VerifierID, v.Notes,
	)
	return err
}

// UpdateVerificationStatus - approve or reject a verification
func UpdateVerificationStatus(db *sql.DB, id string, status string, notes string) error {
	_, err := db.Exec(
		"UPDATE verifications SET status = $1, notes = $2 WHERE id = $3",
		status, notes, id,
	)
	return err
}