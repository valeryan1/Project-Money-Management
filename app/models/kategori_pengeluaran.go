package models

import (
	"database/sql"
)

type KategoriPengeluaran struct {
	KategoriPengeluaranID int
	Kategori              string
}

func (kp *KategoriPengeluaran) Create(db *sql.DB) error {
	query := `INSERT INTO KategoriPengeluaran (Kategori) VALUES (?)`
	_, err := db.Exec(query, kp.Kategori)
	return err
}

func GetKategoriPengeluarans(db *sql.DB) ([]KategoriPengeluaran, error) {
	query := `SELECT KategoriPengeluaranID, Kategori FROM KategoriPengeluaran`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []KategoriPengeluaran
	for rows.Next() {
		var category KategoriPengeluaran
		if err := rows.Scan(&category.KategoriPengeluaranID, &category.Kategori); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
