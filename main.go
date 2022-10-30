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
		w.WriteHeader(400)
		return
		//panic("Cannot Find User")
	}

	// Createing multiple rows for each user with single tag
	for _, tag := range req.Tags {
		var uTag UserTag
		uTag.User, uTag.Tag = user, tag
		DB.Create(&uTag)

		var uTot UserTot
		uTot.User, uTot.UserTag, uTot.Name = user, uTag, user.FirstName+" "+user.LastName
		timeNowInMilli := (time.Now().UnixNano() / int64(time.Millisecond))
		uTot.ExpiryTime = timeNowInMilli + (req.Expiry * int64(time.Millisecond)) // Expiry Time in Milliseconds
		DB.Create(&uTot)
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
