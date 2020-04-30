package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type user struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Age      int    `json:"Age"`
	Job      string `json:"Job"`
	Friendly bool   `json:"Friendly"`
}

type allUsers []user

var u = allUsers{}

//connects to the databa
func connectDB() {
	db, _ = sql.Open("postgres", "postgres://postgres:someone@localhost/api-test?sslmode=disable")
}

//initializes the data from the database into the user struct
func initialize(rows *sql.Rows) {
	var temp user
	var empty allUsers
	u = empty
	i := 0
	for rows.Next() {
		u = append(u, temp)
		rows.Scan(&u[i].ID, &u[i].Name, &u[i].Age, &u[i].Job, &u[i].Friendly)

		i = i + 1
	}

}

//returns all the data in the database
func viewAll(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		fmt.Println("Couldn't make it into rows")
		panic(err)
	}
	initialize(rows)
	for i, _ := range u {
		fmt.Fprintln(w, u[i].ID, u[i].Name, u[i].Age, u[i].Job, u[i].Friendly)
	}

}

//creates a user data
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser user

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	_ = body
	json.Unmarshal(body, &newUser)
	query := "Insert into users values('" + newUser.ID + "','" + newUser.Name + "'," + strconv.Itoa(newUser.Age) + ",'" + newUser.Job + "'," + "true" + ")"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Couldn't make it into rows")
		panic(err)
	}

	initialize(rows)
	for i, _ := range u {
		fmt.Fprintln(w, u[i].ID, u[i].Name, u[i].Age, u[i].Job, u[i].Friendly)
	}

}

//returns the user data based on ID
func getOneUser(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["ID"]

	rows, err := db.Query(fmt.Sprintf("Select * from users where id = %v", userID))
	if err != nil {
		fmt.Println("Couldn't make it into rows")
		panic(err)
	}
	initialize(rows)
	for i, _ := range u {
		fmt.Fprintln(w, u[i].ID, u[i].Name, u[i].Age, u[i].Job, u[i].Friendly)
	}

}

//updates the user data
func updateUser(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["ID"]
	var newUser user

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	_ = body
	json.Unmarshal(body, &newUser)
	rows, err := db.Query(fmt.Sprintf("Update  users set age = %v where ID=%v", newUser.Age, userID))
	if err != nil {
		fmt.Println("Couldn't make it into rows")
		panic(err)
	}
	_ = rows

}

//deletes a user data
func deleteUser(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["ID"]

	rows, err := db.Query(fmt.Sprintf("Delete from users  where ID=%v", userID))
	if err != nil {
		fmt.Println("Couldn't make it into rows")
		panic(err)
	}
	_ = rows

}

func main() {

	connectDB()

	router := mux.NewRouter()
	router.HandleFunc("/", viewAll)
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/user/{ID}", getOneUser).Methods("GET")
	router.HandleFunc("/user/{ID}", updateUser).Methods("PATCH")
	router.HandleFunc("/user/{ID}", deleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8001", router))
	defer db.Close()
}
