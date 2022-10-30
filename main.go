package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

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

type UserTag struct {
	gorm.Model
	UserId     uint   `json:userId`
	Name       string `json:"fullName"`
	Tag        string `json:"tag"`
	ExpiryTime string `json:"expiryTime"`
}

type GetResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"fullName"`
	Phone string `json:"phone"`
}

type PostResponse struct {
	ID uint `json:"id"`
}

type TagPostReq struct {
	Tags   []string `json:"tags"`
	Expiry uint     `json:"expiry"`
}

func intialMigration() {
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("cannot connect DB")
	}
	DB.AutoMigrate(&User{}, &UserTag{})
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	var users []User
	DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewUser")
	w.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	DB.Create(&user)
	fmt.Println(user)
	var res PostResponse
	res.ID = user.ID
	json.NewEncoder(w).Encode(res)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateUser")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	DB.First(&user, params["id"])
	json.NewDecoder(r.Body).Decode(&user)
	DB.Save(&user)
	fmt.Println(user)
	json.NewEncoder(w).Encode(user)
}

func updateUserTag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateUserTag")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var req TagPostReq
	json.NewDecoder(r.Body).Decode(&req) // Decoding tags & expiry

	var user User
	result := DB.First(&user, params["id"])
	fmt.Println(result)
	if result.Error != nil {
		fmt.Println(result.Error)
		w.WriteHeader(404)
		return
		//panic("Cannot Find User")
	}

	// Createing multiple rows for each user with single tag
	for ind, tag := range req.Tags {
		var uTag []UserTag
		uTag.UserId, uTag.Name = user.ID, user.FirstName+" "+user.LastName
		uTag.Tag = tag

		// Setting an Expiry time by adding the duration with reqest submit time, later it will be compared
		uTag.ExpiryTime = time.Now().Add(time.Duration(req.Expiry) * time.Millisecond)
		fmt.Println(time.Now())
		fmt.Println(time.Now().Add(time.Duration(req.Expiry) * time.Millisecond))
		DB.Create(&uTag)
		uTag.ID = user.ID
		DB.Save(&uTag)
	}

}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getUser")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	result := DB.First(&user, params["id"])
	fmt.Println(result)
	if result.Error != nil {
		fmt.Println(result.Error)
		w.WriteHeader(404)
		return
		//panic("Cannot Find User")
	}
	var res GetResponse
	res.ID, res.Name, res.Phone = user.ID, user.FirstName+" "+user.LastName, user.Phone
	json.NewEncoder(w).Encode(res)
}

func errorHandler(w http.ResponseWriter, r *http.Request, i int) {
	panic("unimplemented")
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
	myRouter.HandleFunc("/users/{id}", getUser).Methods("GET")

	myRouter.HandleFunc("/users", returnAllUsers).Methods("GET")
	myRouter.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/users", updateUser).Methods("PUT")
	http.Handle("/", myRouter)

	// Extended Tasks
	myRouter.HandleFunc("/users/{id}/tags", updateUserTag).Methods("POST")

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
