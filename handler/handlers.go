package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"myapp/model"
	"net/http"
	"strconv"

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
		ID SERIAL PRIMARY KEY,
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
	a.Router.Path("/search").Queries("name", "{name}").HandlerFunc(a.SearchUserByName).Methods("GET")
	a.Router.HandleFunc("/delete/{id}", a.DeleteUser).Methods("GET")
	a.Router.HandleFunc("/modify/{id}", a.ModifyUser).Methods("POST")
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

func (a *App) GenerateID() int64 {
	query := "SELECT COUNT (*) FROM users"
	res, _ := a.DB.Exec(query)
	ID, _ := res.LastInsertId()
	fmt.Print(res)
	return ID
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

	id, err := model.CreateUser(a.DB, u)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.ID = id
	RespondJSON(w, http.StatusCreated, u)
}

func (a *App) SearchUserByName(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	name := params["name"]
	users, err := model.SearchUserByName(a.DB, name)

	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, users)
}

func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := model.DeleteUser(a.DB, id); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"message": "User deleted"})
}

func (a *App) ModifyUser(w http.ResponseWriter, r *http.Request) {
	var u model.Users
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.ID = id

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()
	if err := model.ModifyUser(a.DB, u); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"message": "User modified"})
}
