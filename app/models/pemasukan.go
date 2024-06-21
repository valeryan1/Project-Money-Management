package models

import (
	"database/sql"
)

type Pemasukan struct {
	PemasukanID int
	UserID      int
	UmkmID      int
	Nominal     float64
	Date        string
}

func (p *Pemasukan) Create(db *sql.DB) error {
	query := `INSERT INTO Pemasukan (UserID, UmkmID, Nominal, Date)
              VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, p.UserID, p.UmkmID, p.Nominal, p.Date)
	return err
}

func GetPemasukans(db *sql.DB) ([]Pemasukan, error) {
	query := `SELECT PemasukanID, UserID, UmkmID, Nominal, Date FROM Pemasukan`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pemasukan []Pemasukan
	for rows.Next() {
		var p Pemasukan
		if err := rows.Scan(&p.PemasukanID, &p.UserID, &p.UmkmID, &p.Nominal, &p.Date); err != nil {
			return nil, err
		}
		pemasukan = append(pemasukan, p)
	}
	return pemasukan, nil
}
