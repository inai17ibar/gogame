module github.com/inai17ibar/gogame

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/mux v1.8.0
)

replace github.com/inai17ibar/gogame/model => ./model
replace github.com/inai17ibar/gogame/tool => ./tool