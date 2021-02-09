package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

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

func createToken(user User) (string, error) {
	var err error

	secret := "secret"

	//jst structure {base64 encoded header. paylead. signature}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Name,
		"iss":  "__init__", //JWTの発行者をいれるべき？
	})

	var tokenString string
	tokenString, err = token.SignedString([]byte(secret))

	fmt.Println("tokenString:", tokenString)

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
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
	var user User
	var error Error

	json.NewDecoder(r.Body).Decode(&user)

	if user.Name == "" {
		error.Message = "Username is invalid."
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is invalid."
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	password := user.Password
	fmt.Println("input password: ", password)

	var time time.Time
	//get auth key of user info from db.
	row := db.QueryRow("select * from gogame_db.user_table where username = ?", &user.Name)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.SessionId, &time)

	if err != nil {
		if err == sql.ErrNoRows { //https://golang.org/pkg/database/sql/#pkg-variables
			error.Message = "Not found User."
			errorInResponse(w, http.StatusBadRequest, error)
		} else {
			log.Fatal(err)
		}
	}

	hasedPassword := user.Password
	fmt.Println("hasedPassword(db): ", hasedPassword)

	//password check
	err = bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))

	if err != nil {
		error.Message = "Invalid Password."
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	token, err := createToken(user)

	if err != nil {
		log.Fatal(err)
	}

	var jwt JWT
	w.WriteHeader(http.StatusOK)
	jwt.Token = token
	responseByJSON(w, jwt)
}

func handleRequests() {
	// ルーターのイニシャライズ
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")
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
