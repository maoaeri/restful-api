package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"myapp/model"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) InitDB(user, password, name string) {
	DBinfo := fmt.Sprintf("host='localhost' port='5432' user=%s password=%s dbname=%s sslmode=disable", user, password, name)
	DB, err := sql.Open("postgres", DBinfo)
	a.DB = DB
	if err != nil {
		panic(err)
	}
	q := `CREATE TABLE IF NOT EXISTS users (
		Name VARCHAR(255),
		Birthday VARCHAR(255),
		Sex VARCHAR(255),
		Email VARCHAR(255)
	);`
	a.DB.Exec(q)
}

func (a *App) InitRouter() {
	a.Router.HandleFunc("/users", a.GetAllUsers).Methods("GET")
	a.Router.HandleFunc("/create/user", a.CreateUser).Methods("POST")
	a.Router.HandleFunc("/foo", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello, %q", html.EscapeString(r.URL.Path))

	}).Methods("GET")
}

func (a *App) Init(user, password, name string) {
	a.InitDB(user, password, name)
	a.InitRouter()
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func (a *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := model.GetAllUsers(a.DB)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}
	RespondJSON(w, http.StatusOK, users)
}

func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u model.Users
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := model.CreateUser(a.DB, u); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, u)
}