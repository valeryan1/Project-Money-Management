package models

import "database/sql"

type Reminder struct {
	ReminderID  int
	UserID      int
	Description string
	DueDate     string
}

func GetReminders(db *sql.DB, userID int) ([]Reminder, error) {
	query := `SELECT ReminderID, UserID, Description, DueDate FROM Reminders WHERE UserID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var reminder Reminder
		if err := rows.Scan(&reminder.ReminderID, &reminder.UserID, &reminder.Description, &reminder.DueDate); err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	return reminders, nil
}
