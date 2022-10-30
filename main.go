package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *sql.DB
var alreadyInitialized bool
var errorInDB bool

// DB connection params
const (
	username = "root"
	password = "mypass98"
	hostname = "127.0.0.1:3306"
	dbname   = "users_db"
)

var DB *gorm.DB
var err error

const DSN = "root:mypass98@tcp(127.0.0.1:3306)/users_db?parseTime=true"

type User struct {
	gorm.Model
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

func intialMigration() {
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("cannot connect DB")
	}
	DB.AutoMigrate(&User{})
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	var userss []User
	DB.Find(&userss)
	json.NewEncoder(w).Encode(userss)

}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewUser")
	w.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	DB.Create(&user)
	fmt.Println(user)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateUser")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.First(&user, params["id"])
	json.NewDecoder(r.Body).Decode(&user)
	DB.Save(&user)
	fmt.Println(user)
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getUser")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.First(&user, params["id"])
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteUser")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.Delete(&user, params["id"])
	json.NewEncoder(w).Encode(user)
}

func handleRequests() {
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/users", createNewUser).Methods("POST")
	myRouter.HandleFunc("/users", updateUser).Methods("PUT")

	myRouter.HandleFunc("/users", returnAllUsers).Methods("GET")
	myRouter.HandleFunc("/users/{id}", getUser).Methods("GET")
	myRouter.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	//myRouter.HandleFunc("/article/{id}", returnSingleArticle)

	http.ListenAndServe(":8000", myRouter)
}

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)

}

func main() {
	intialMigration()
	handleRequests()
}
