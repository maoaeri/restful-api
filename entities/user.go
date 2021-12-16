package entities

import (
	"database/sql"
	"fmt"
)

type User struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Sex      string `json:"sex"`
	Email    string `json:"email"`
}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "vvlalalove193"
	DB_NAME     = "users"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, _ := sql.Open("postgres", dbinfo)

	return db
}

func main() {
	setupDB()
}
