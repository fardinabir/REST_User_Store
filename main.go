package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// User ...
type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	Id        int    `json:"id"`
}

type NewReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

// Users slice
var Users []User

func handleRequests() {
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/users", createNewUser2).Methods("POST")
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", returnAllUsers).Methods("GET")

	myRouter.HandleFunc("/usersf/{id}", deleteUser).Methods("DELETE")
	//myRouter.HandleFunc("/article/{id}", returnSingleArticle)

	http.ListenAndServe(":8000", myRouter)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Users)
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewUsers")
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(reqBody)
	var newUserReq NewReq
	json.Unmarshal(reqBody, &newUserReq)
	Users = append(Users, assignIdNewUser(newUserReq))

	//	fmt.Println(Users)
	json.NewEncoder(w).Encode(Users[len(Users)-1])
}

func createNewUser2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usr User
	_ = json.NewDecoder(r.Body).Decode(&usr)
	fmt.Println(usr)
	usr.Id = rand.Intn(1000000)
	Users = append(Users, usr)
	json.NewEncoder(w).Encode(usr)
}

func assignIdNewUser(nr NewReq) User {
	var usr User
	usr.FirstName, usr.LastName, usr.Password, usr.Phone = nr.FirstName, nr.LastName, nr.Password, nr.Password
	usr.Id = len(Users)
	return usr
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteArticle")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for index, user := range Users {
		if user.Id == id {
			Users = append(Users[:index], Users[index+1:]...)
		}
	}

}

func main() {
	Users = []User{
		{
			FirstName: "fardin",
			LastName:  "abir",
			Password:  "1234",
			Phone:     "131424243",
		},
	}

	handleRequests()
}
