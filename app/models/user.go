package models

import (
	"database/sql"
	"time"
)

type User struct {
	UserID       int
	UmkmID       int
	Name         string
	Username     string
	UserAddress  string
	UserEmail    string
	UserPhone    string
	UserDOB      time.Time
	UserPassword string
}

func (u *User) Create(db *sql.DB) error {
	query := `INSERT INTO user (Name, Username, UserAddress, UserEmail, UserPhone, UserDOB, UserPassword)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, u.Name, u.Username, u.UserAddress, u.UserEmail, u.UserPhone, u.UserDOB, u.UserPassword)
	return err
}

// func GetUserByUsername(db *sql.DB, username string) (*User, error) {
// 	var user User
// 	query := `SELECT UserID, UmkmID, Name, Username, UserAddress, UserEmail, UserPhone, UserDOB, UserPassword
//               FROM user WHERE Username = ?`
// 	err := db.QueryRow(query, username).Scan(&user.UserID, &user.UmkmID, &user.Name, &user.Username, &user.UserAddress, &user.UserEmail, &user.UserPhone, &user.UserDOB, &user.UserPassword)
// 	return &user, err
// }

// GetUserByUsername queries the database for a user with the given username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT UserID, UmkmID, Name, Username, UserAddress, UserEmail, UserPhone, UserDOB, UserPassword FROM user WHERE Username = ?", username).Scan(
		&user.UserID, &user.UmkmID, &user.Name, &user.Username, &user.UserAddress, &user.UserEmail, &user.UserPhone, &user.UserDOB, &user.UserPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil user and nil error if no rows found
		}
		return nil, err
	}
	return &user, nil
}
