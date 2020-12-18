package main

import (
	"fmt"
	"log"
	"net/http"
)

type UserInfo struct {
	Name     string `json:"name"`
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call signup.")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call login")
}

// func createUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("post CreatUser request.")
// }

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("get User request.")
// 	json.NewEncoder(w).Encode(UserInfo.Name)
// }

// func updateUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("put UpdateUser request.")
// }

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Api demo!")
	//sign up and sign in page
	fmt.Println("Endpoint Hit: demo page")
}

func handleRequests() {
	http.HandleFunc("/", home)

	//webぺージにリダイレクトする
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)

	log.Println("start api server :received 8080 port")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
