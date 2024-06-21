package models

import (
	"database/sql"
)

type Umkm struct {
	UmkmID         int
	UmkmName       string
	UmkmCategoryID int
	UmkmAddress    string
	UmkmNoTelp     string
}

func (u *Umkm) Create(db *sql.DB) error {
	query := `INSERT INTO UMKM (UmkmName, UmkmCategoryID, UmkmAddress, UmkmNoTelp)
              VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, u.UmkmName, u.UmkmCategoryID, u.UmkmAddress, u.UmkmNoTelp)
	return err
}

func GetUmkms(db *sql.DB) ([]Umkm, error) {
	query := `SELECT UmkmID, UmkmName, UmkmCategoryID, UmkmAddress, UmkmNoTelp FROM UMKM`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var umkms []Umkm
	for rows.Next() {
		var umkm Umkm
		if err := rows.Scan(&umkm.UmkmID, &umkm.UmkmName, &umkm.UmkmCategoryID, &umkm.UmkmAddress, &umkm.UmkmNoTelp); err != nil {
			return nil, err
		}
		umkms = append(umkms, umkm)
	}
	return umkms, nil
}

func GetUMKMsByUserID(db *sql.DB, userID int) ([]Umkm, error) {
	query := `SELECT UmkmID, UmkmName FROM UMKM WHERE UserID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var umkms []Umkm
	for rows.Next() {
		var umkm Umkm
		if err := rows.Scan(&umkm.UmkmID, &umkm.UmkmName); err != nil {
			return nil, err
		}
		umkms = append(umkms, umkm)
	}
	return umkms, nil
}
