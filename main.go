package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"./tool"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var err error

const (
	// secret は openssl rand -base64 40 コマンドで作成
	secret = "secret" //電子署名作成のためのシークレット、外部にもれないようにしたい

	// aud をつかっているが、独自のキーを作ったほうがいい？
	userKey = "aud"

	// iat と exp は登録済みクレーム名。それぞれの意味は https://tools.ietf.org/html/rfc7519#section-4.1 を参照。{
	iatKey = "iat"
	expKey = "exp"

	// lifetime は jwt の発行から失効までの期間を表す。
	lifetime = 30 * time.Minute
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

func createToken(userName string) (string, error) {

	//jwt structure {base64 encoded header. paylead. signature}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":   "AccessToken",
			"iss":   "https://idp.gogame.com/", // 発行者
			userKey: userName,
			iatKey:  time.Now().Unix(),
			expKey:  time.Now().Add(time.Hour * 72).Unix(),
		})

	//電子署名
	var tokenString string
	tokenString, err = token.SignedString([]byte(secret))

	//fmt.Println("tokenString:", tokenString)

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := errors.New("Unexpected signing method")
			return nil, err
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, errors.New("Token is invalid")
	}
	return parsedToken, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home ")
}

func createUserInfo(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	var jwt JWT

	json.NewDecoder(r.Body).Decode(&user)

	if user.Name == "" {
		error.Message = "Nameは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	// if user.Password == "" {
	// 	error.Message = "パスワードは必須です。"
	// 	errorInResponse(w, http.StatusBadRequest, error)
	// 	return
	// }

	// hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// user.Password = string(hash)

	//regist database
	//sql_query := "INSERT INTO gogame_db.user_table (username, password) VALUES (?, ?)"
	sql_query := "INSERT INTO gogame_db.user_table (username) VALUES (?)"
	_, err = db.Exec(sql_query, user.Name)

	if err != nil {
		//usernameがすでにあった場合

		//その他のエラー
		fmt.Println(err)
		error.Message = "SQLServer Error"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//generate token
	token, err := createToken(user.Name)
	if err != nil {
		log.Fatal(err)
	}

	//DBに登録後パスワードは空にする
	//user.Password = ""

	jwt.Token = token
	w.WriteHeader(http.StatusOK)
	responseByJSON(w, jwt)
}

func verify(w http.ResponseWriter, r *http.Request) {
	var token *jwt.Token

	//headerからtokenをとりだし、検証する
	token, err := verifyToken(r.Header.Get("x-token"))
	if err != nil {
		//errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	fmt.Println(token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("not found claims")
		return
	}
	fmt.Println(claims)
	iat, ok := claims[iatKey].(float64)
	if !ok {
		fmt.Println("not found claims of iat")
		return
	}
	fmt.Println(iat)
	userID, ok := claims[userKey].(string)
	if !ok {
		fmt.Println("not found claims of userid")
		return
	}
	fmt.Println(userID)

	w.Header().Set("Content-Type", "application/json")
	responseByJSON(w, nil)
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	var token *jwt.Token

	//headerからtokenをとりだし、検証する
	token, err := verifyToken(r.Header.Get("x-token"))

	//get userinfo from db.
	row := db.QueryRow("select * from gogame_db.user_table where username = ?", &token.Claims)
	err = row.Scan(&user.ID, &user.Name, &user.Password, &user.Created_date)

	if err != nil {
		if err == sql.ErrNoRows { //https://golang.org/pkg/database/sql/#pkg-variables
			error.Message = "Not found User."
			errorInResponse(w, http.StatusBadRequest, error)
		} else {
			log.Fatal(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	responseByJSON(w, user)
}

func handleRequests() {
	// ルーターのイニシャライズ
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/user/create", createUserInfo).Methods("POST")
	router.HandleFunc("/user/get", getUserInfo).Methods("GET")
	router.HandleFunc("/verify", verify).Methods("GET")

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
