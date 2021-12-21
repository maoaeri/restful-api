package model

import (
	"database/sql"
	"fmt"
)

type Users struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Sex      string `json:"sex"`
	Email    string `json:"email"`
}

func GetAllUsers(db *sql.DB) ([]Users, error) {
	rows, err := db.Query("SELECT * FROM users;")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []Users{}

	for rows.Next() {
		var u Users
		if err := rows.Scan(&u.Name, &u.Birthday, &u.Sex, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func CreateUser(db *sql.DB, u Users) error {
	query := fmt.Sprintf("INSERT INTO users VALUE (%v, %v, %v, %v )", u.Name, u.Birthday, u.Sex, u.Birthday)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
