package model

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/inai17ibar/gogame/tool"

	_ "github.com/go-sql-driver/mysql"
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

// dbインスタンス格納用
var db *sql.DB

func selectWeightRandom(itemWeights []int, times int) (results []int) {
	sort.Ints(itemWeights) // 重みの昇順でチェックしていくので重み一覧をソートする

	// 抽選結果チェックの基準となる境界値を生成
	boundaries := make([]int, len(itemWeights)+1)
	for i := 1; i < len(itemWeights)+1; i++ {
		boundaries[i] = boundaries[i-1] + itemWeights[i-1]
	}
	boundaries = boundaries[1:len(boundaries)] // 頭の0値は不要なので破棄

	rand.Seed(time.Now().UnixNano())
	// 境界検知用のスライス
	result := make([]int, len(itemWeights))
	n := 1
	for i := 0; i < n; i++ {
		draw := rand.Intn(boundaries[len(boundaries)-1]) + 1
		for i, boundary := range boundaries {
			if draw <= boundary {
				result[i]++
				results = append(results, itemWeights[i])
				break
			}
		}
	}

	// for i := 0; i < len(itemWeights); i++ {
	// 	fmt.Printf("%4d  %f  %f\n",
	// 		itemWeights[i],
	// 		float64(itemWeights[i])/float64(boundaries[len(boundaries)-1]),
	// 		float64(result[i])/float64(n),
	// 	)
	// }
	return results
}

func ExecuteGacha(times int) ([]GachaResult, error) {
	var results []GachaResult
	var err error

	// parmas.go から DB の URL を取得
	info := tool.DatabaseInfo{}
	db, err = sql.Open("mysql", info.GetDBUrl())
	if err != nil {
		panic(err)
	}
	log.Println("Connected to mysql.")

	var id int
	var name string
	var gtype int
	var start, end time.Time

	weightRarity := []int{0, 0, 0, 0, 0}
	//DBからレアリティ抽選の重みを読み込む
	row, err := db.Query("select * from gogame_db.gacha_table where `name` = 'sample_gacha'")
	if err != nil {
		panic(err)
	}
	defer row.Close()
	// err = row.Scan(&id, &name, &gtype, &weightRarity[0],
	// 	&weightRarity[1], &weightRarity[2], &weightRarity[3], &weightRarity[4], &start, &end)
	for row.Next() {
		err := row.Scan(&id, &name, &gtype, &weightRarity[0],
			&weightRarity[1], &weightRarity[2], &weightRarity[3], &weightRarity[4], &start, &end)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, name)
	}

	if err != nil {
		if err == sql.ErrNoRows { //https://golang.org/pkg/database/sql/#pkg-variables
			log.Panicln("Not found data-row in database.")
		} else {
			log.Fatal(err)
		}
	}

	fmt.Println(weightRarity[0], weightRarity[1])

	return results, nil

	//回数だけ抽選を実行する
	for i := 0; i < times; i++ {

	}

	//DBからキャラ抽選の重みとデータを読み込む
	_, err = db.Exec("select from gogame_db.gacha_characters_01_table")

	if err != nil {
		if err == sql.ErrNoRows { //https://golang.org/pkg/database/sql/#pkg-variables
			log.Panicln("Not found data-row in database.")
		} else {
			log.Fatal(err)
		}
	}

	//回数だけ抽選を実行する
	for i := 0; i < times; i++ {

	}

	return results, nil
}
