package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home ")
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup.")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login")
}

func handleRequests() {
	// ルーターのイニシャライズ
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/signup", signupHandler)
	router.HandleFunc("/login", loginHandler)
	http.Handle("/", router)

	log.Println("start api server :received 8080 port")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
