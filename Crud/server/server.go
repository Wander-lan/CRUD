package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Create a new user to put on database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error at reading request body"))
		return
	}

	var user user

	if err = json.Unmarshal(requestBody, &user); err != nil {
		w.Write([]byte("Error at converting user to struct"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.Write([]byte("Error at connecting with database"))
	}
	defer db.Close()

	// Prepare statement:
	statement, err := db.Prepare("insert into users (name, email) values (?, ?)")
	if err != nil {
		w.Write([]byte("Error at creating statement"))
	}
	defer statement.Close()

	insertion, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.Write([]byte("Error at executing the statement"))
		return
	}

	idInserted, err := insertion.LastInsertId()
	if err != nil {
		w.Write([]byte("Error at getting the inserted Id"))
		return
	}

	// Status codes
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User inserted! Id: %d", idInserted)))
}

// Find users (plural) from the database
func FindUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		w.Write([]byte("Error at connecting with database"))
	}
	defer db.Close()

	query, err := db.Query("select * from users")
	if err != nil {
		w.Write([]byte("Error at finding users"))
	}
	defer query.Close()

	var users []user
	for query.Next() {
		var user user

		if err := query.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Error at scanning the user"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Error at converting users to JSON"))
		return
	}
}

// Find a specific user from the database
func FindUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error at converting parameter to int"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.Write([]byte("Error at connecting with the database"))
		return
	}

	query, err := db.Query("select * from users where id = ?", ID)
	if err != nil {
		w.Write([]byte("Error at finding the user"))
		return
	}

	var user user
	if query.Next() {
		if err := query.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Error at scanning the user"))
			return
		}
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Error at converting user to JSON"))
		return
	}
}