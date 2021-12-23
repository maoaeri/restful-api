package model

import (
	"database/sql"
	"fmt"
)

type Users struct {
	ID       int    `json:"ID"`
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
		if err := rows.Scan(&u.ID, &u.Name, &u.Birthday, &u.Sex, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func CreateUser(db *sql.DB, u Users) (id int, err error) {
	query := fmt.Sprintf("INSERT INTO users VALUES (DEFAULT, '%v', '%v', '%v', '%v') RETURNING ID", u.Name, u.Birthday, u.Sex, u.Email)
	err = db.QueryRow(query).Scan(&id)
	if err != nil {
		return id, err
	}
	//ID, err := res.LastInsertId()
	//u.ID = ID
	fmt.Printf("dm")
	return id, nil
}

func DeleteUser(db *sql.DB, id int) error {
	query := fmt.Sprintf("DELETE FROM users WHERE id=%v", id)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func SearchUserByName(db *sql.DB, name string) ([]Users, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE name='%v';", name)
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []Users{}

	for rows.Next() {
		var u Users
		if err := rows.Scan(&u.ID, &u.Name, &u.Birthday, &u.Sex, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func ModifyUser(db *sql.DB, u Users) error {
	query := fmt.Sprintf("UPDATE users SET name='%v', birthday='%v', sex='%v', email='%v' WHERE ID=%v", u.Name, u.Birthday, u.Sex, u.Email, u.ID)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
