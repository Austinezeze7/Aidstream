package models

import (
	"database/sql"
	"time"
)

type Donation struct {
	ID              string    `json:"id"`
	DonorID         string    `json:"donor_id"`
	CaseID          string    `json:"case_id"`
	Amount          float64   `json:"amount"`
	TransactionHash string    `json:"transaction_hash"`
	Status          string    `json:"status"`
	DonatedAt       time.Time `json:"donated_at"`
}

// GetAllDonations - fetch all donations
func GetAllDonations(db *sql.DB) ([]Donation, error) {
	rows, err := db.Query("SELECT id, donor_id, case_id, amount, transaction_hash, status, donated_at FROM donations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		err := rows.Scan(&d.ID, &d.DonorID, &d.CaseID, &d.Amount, &d.TransactionHash, &d.Status, &d.DonatedAt)
		if err != nil {
			return nil, err
		}
		donations = append(donations, d)
	}
	return donations, nil
}

// GetDonationsByCase - fetch all donations for a specific case
func GetDonationsByCase(db *sql.DB, caseID string) ([]Donation, error) {
	rows, err := db.Query("SELECT id, donor_id, case_id, amount, transaction_hash, status, donated_at FROM donations WHERE case_id = $1", caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		err := rows.Scan(&d.ID, &d.DonorID, &d.CaseID, &d.Amount, &d.TransactionHash, &d.Status, &d.DonatedAt)
		if err != nil {
			return nil, err
		}
		donations = append(donations, d)
	}
	return donations, nil
}

// CreateDonation - insert a new donation
func CreateDonation(db *sql.DB, d Donation) error {
	_, err := db.Exec(
		"INSERT INTO donations (donor_id, case_id, amount, transaction_hash) VALUES ($1, $2, $3, $4)",
		d.DonorID, d.CaseID, d.Amount, d.TransactionHash,
	)
	return err
}

// UpdateDonationStatus - update donation status
func UpdateDonationStatus(db *sql.DB, id string, status string) error {
	_, err := db.Exec("UPDATE donations SET status = $1 WHERE id = $2", status, id)
	return err
}