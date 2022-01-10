package main

import (
	"myapp/handler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "123456"
	DB_NAME     = "postgres"
)

func main() {
	app := &handler.App{}
	app.Router = mux.NewRouter()
	app.Init(DB_USER, DB_PASSWORD, DB_NAME)
	app.Run()
}
