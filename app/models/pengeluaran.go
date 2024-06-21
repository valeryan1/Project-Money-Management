package models

import (
	"database/sql"
)

// Pengeluaran represents an expense record
type Pengeluaran struct {
	PengeluaranID         int
	UserID                int
	UmkmID                int
	Nominal               float64
	Date                  string
	KategoriPengeluaranID int
}

// GetPengeluarans retrieves the expenses for a specific user
func GetPengeluarans(db *sql.DB, userID int) ([]Pengeluaran, error) {
	query := `SELECT PengeluaranID, UserID, UmkmID, Nominal, Date, KategoriPengeluaranID FROM Pengeluaran WHERE UserID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pengeluarans []Pengeluaran
	for rows.Next() {
		var pengeluaran Pengeluaran

		// Declare dateStr variable to store Date as string
		// var dateStr string
		if err := rows.Scan(&pengeluaran.PengeluaranID, &pengeluaran.UserID, &pengeluaran.UmkmID, &pengeluaran.Nominal, &pengeluaran.Date, &pengeluaran.KategoriPengeluaranID); err != nil {
			return nil, err
		}
		// Parsing dateStr into time.Time
		// pengeluaran.Date, err = time.Parse("2006-01-02", dateStr)

		pengeluarans = append(pengeluarans, pengeluaran)
	}
	return pengeluarans, nil
}
