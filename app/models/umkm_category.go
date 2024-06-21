package models

import (
	"database/sql"
)

type UmkmCategory struct {
	UmkmCategoryID int
	Category       string
}

func (uc *UmkmCategory) Create(db *sql.DB) error {
	query := `INSERT INTO UMKMCategory (Category) VALUES (?)`
	_, err := db.Exec(query, uc.Category)
	return err
}

func GetUmkmCategories(db *sql.DB) ([]UmkmCategory, error) {
	query := `SELECT UmkmCategoryID, Category FROM UMKMCategory`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []UmkmCategory
	for rows.Next() {
		var category UmkmCategory
		if err := rows.Scan(&category.UmkmCategoryID, &category.Category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
