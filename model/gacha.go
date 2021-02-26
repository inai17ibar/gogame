package model

import (
	"database/sql"
	"log"

	"github.com/inai17ibar/gogame/tool"
)

type GachaResult struct {
	CharacterId string `json:"character_id"`
	Name        string `json:"name"`
}

type DrawGachaRequest struct {
	Times int `json:"times"`
}

type DrawGachaResponse struct {
	Results []GachaResult `json:"times"`
}

type Gacha struct {
}

// dbインスタンス格納用
var db *sql.DB

func (g Gacha) executeGacha(times int) ([]GachaResult, error) {
	var results []GachaResult
	var err error

	// parmas.go から DB の URL を取得
	info := tool.DatabaseInfo{}
	db, err = sql.Open("mysql", info.GetDBUrl())
	if err != nil {
		panic(err)
	}

	//回数だけselectを実行する
	_, err = db.Exec("select from gogame_db.gacha_characters_01_table")

	if err != nil {
		if err == sql.ErrNoRows { //https://golang.org/pkg/database/sql/#pkg-variables
			log.Panicln("Not found data-row in database.")
		} else {
			log.Fatal(err)
		}
	}

	return results, nil
}
