package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	. "github.com/inai17ibar/gogame/model"
	"github.com/inai17ibar/gogame/tool"

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

func createToken(userId string) (string, error) {

	//jwt structure {base64 encoded header. paylead. signature}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":   "AccessToken",
			"iss":   "https://idp.gogame.com/", // 発行者
			userKey: userId,
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
		// check signing method　型アサーションを使ってアルゴリズムをチェックしている
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := errors.New("Unexpected signing method")
			return nil, err
		}
		return []byte("secret"), nil //検証に使うキーをいれる
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

	//useridの取得
	row := db.QueryRow("SELECT LAST_INSERT_ID()")
	err = row.Scan(&user.ID)
	if err != nil {
		log.Fatal(err)
	}

	//generate token
	token, err := createToken(strconv.Itoa(user.ID))
	if err != nil {
		log.Fatal(err)
	}

	//DBに登録後パスワードは空にする
	//user.Password = ""

	jwt.Token = token
	w.WriteHeader(http.StatusOK)
	responseByJSON(w, jwt)
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error

	tokenString := r.Header.Get("x-token")
	if tokenString == "" {
		fmt.Println("Not found token")
		return
	}

	//headerからtokenをとりだし、検証する
	token, err := verifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	userID, ok := claims[userKey].(string)
	if !ok {
		fmt.Println("not found claims or userid")
	}

	//get userinfo from db.
	row := db.QueryRow("select * from gogame_db.user_table where id = ?", userID)
	err = row.Scan(&user.ID, &user.Name, &user.Password, &user.Created_date)
	user.Password = "" // Passwordはいまはつかっていないが一応空にしておく
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

func updateUserInfo(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error

	json.NewDecoder(r.Body).Decode(&user)
	if user.Name == "" {
		error.Message = "Nameは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	tokenString := r.Header.Get("x-token")
	if tokenString == "" {
		fmt.Println("Not found token")
		return
	}

	//headerからtokenをとりだし、検証する
	token, err := verifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userID, ok := claims[userKey].(string)
	if !ok {
		fmt.Println("not found claims or userid")
	}

	//update values on db.
	_, err = db.Exec("update gogame_db.user_table set username = ? where id = ?", user.Name, userID)

	if err != nil {
		if err == sql.ErrNoRows { //https://golang.org/pkg/database/sql/#pkg-variables
			error.Message = "Not found User."
			errorInResponse(w, http.StatusBadRequest, error)
		} else {
			log.Fatal(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	responseByJSON(w, "")
}

func getCharactersList(w http.ResponseWriter, r *http.Request) {

}

func drawGachaHandler(w http.ResponseWriter, r *http.Request) {
	var request DrawGachaRequest
	var error Error

	json.NewDecoder(r.Body).Decode(&request)
	if request.Times == 0 {
		error.Message = "timesが0なので、ガチャを実行できません"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	fmt.Println("times:", request.Times)

	tokenString := r.Header.Get("x-token")
	if tokenString == "" {
		fmt.Println("Not found token")
		return
	}

	//headerからtokenをとりだし、検証する
	token, err := verifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userID, ok := claims[userKey].(string)
	if !ok {
		fmt.Println("not found claims or userid")
	}
	_ = userID

	fmt.Println("Start gacha")
	ExecuteGacha(request.Times)

	//draw gacha
	// _, err = db.Exec("update gogame_db.user_table set username = ? where id = ?", user.Name, userID)

	// if err != nil {
	// 	if err == sql.ErrNoRows { //https://golang.org/pkg/database/sql/#pkg-variables
	// 		error.Message = "Not found User."
	// 		errorInResponse(w, http.StatusBadRequest, error)
	// 	} else {
	// 		log.Fatal(err)
	// 	}
	// }

	w.Header().Set("Content-Type", "application/json")
	responseByJSON(w, "")
}

func handleRequests() {
	// ルーターのイニシャライズ
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/user/create", createUserInfo).Methods("POST")
	router.HandleFunc("/user/get", getUserInfo).Methods("GET")
	router.HandleFunc("/user/update", updateUserInfo).Methods("PUT")

	router.HandleFunc("/gacha/draw", drawGachaHandler).Methods("POST")
	router.HandleFunc("/character/list", getCharactersList).Methods("GET")

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
