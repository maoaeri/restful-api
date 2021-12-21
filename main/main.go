package main

import (
	"myapp/handler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "vvlalalove193"
	DB_NAME     = "postgres"
)

func main() {
	app := &handler.App{}
	app.Router = mux.NewRouter()
	app.Init(DB_USER, DB_PASSWORD, DB_NAME)
	app.Run()

	/*dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/foo", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello, %q", html.EscapeString(r.URL.Path))

	}).Methods("GET")
	router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	//start and listen to requests
	log.Fatal(http.ListenAndServe(":8080", router))*/
}
