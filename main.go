package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"./tool"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func errorInResponse(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
	return
}

// data interface{} とすると、どのような変数の型でも引数として受け取ることができる
func responseByJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
	return
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home ")
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error

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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password: ", user.Password)
	fmt.Println("Hashed Password", hash)

	user.Password = string(hash)
	fmt.Println("Converted Password:", user.Password)

	info := tool.DatabaseInfo{}
	db, err = sql.Open("mysql", info.GetDBUrl())
	if err != nil {
		panic(err)
	}
	log.Println("Connected to mysql.")

	sql_query := "INSERT INTO gogame_db.user_table (username, password) VALUES (?, ?)"
	_, err = db.Exec(sql_query, user.Name, user.Password)

	if err != nil {
		error.Message = "SQLServer Error"
		fmt.Println(err)
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//DBに登録後パスワードは空にする
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	responseByJSON(w, user)
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
	// info := tool.DatabaseInfo{}

	// var err error
	// // sql.Open("mysql", "user:password@/dbname")
	// db, err = sql.Open("mysql", info.GetDBUrl())
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("Connected to mysql.")

	// //データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	// rows, err := db.Query("SELECT * FROM gogame_db.user_table")

	// if err != nil {
	// 	panic(err.Error())
	// }

	// for rows.Next() {
	// 	user := User{ID: 0, Name: "", Password: "", SessionId: 0}
	// 	var time time.Time
	// 	err := rows.Scan(&user.ID, &user.Name, &user.Password, &user.SessionId, &time)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	fmt.Println(user.ID, user.Name, user.Password, time)
	// }

	handleRequests()
}
