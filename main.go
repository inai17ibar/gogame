package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"./tool"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func errorInResponse(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
	return
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home ")
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error

	fmt.Println(r.Body)

	json.NewDecoder(r.Body).Decode(&user)

	if user.Name == "" {
		error.Message = "Nameは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "パスワードは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	//check data
	fmt.Println(user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login")
}

func handleRequests() {
	// ルーターのイニシャライズ
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler)
	http.Handle("/", router)

	log.Println("start api server :received 8080 port")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// dbインスタンス格納用
var db *sql.DB

func main() {
	// parmas.go から DB の URL を取得
	info := tool.DatabaseInfo{}

	var err error
	// sql.Open("mysql", "user:password@/dbname")
	db, err = sql.Open("mysql", info.GetDBUrl())
	if err != nil {
		panic(err)
	}
	log.Println("Connected to mysql.")

	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	rows, err := db.Query("SELECT * FROM gogame_db.user_table")
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		user := User{ID: 0, Name: "", Password: "", SessionId: 0}
		var time time.Time
		err := rows.Scan(&user.ID, &user.Name, &user.Password, &user.SessionId, &time)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(user.ID, user.Name, user.Password, time)
	}

	handleRequests()
}
