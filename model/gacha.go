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

func selectWeightRandom(itemWeights []int, times int) (lotteryResults []int) {
	sort.Ints(itemWeights) // 重みの昇順でチェックしていくので重み一覧をソートする
	//TODO: 昇順で与えられなかった場合は警告する
	//timesが0でもだめ

	// 抽選結果チェックの基準となる境界値を生成
	boundaries := make([]int, len(itemWeights)+1)
	for i := 1; i < len(itemWeights)+1; i++ {
		boundaries[i] = boundaries[i-1] + itemWeights[i-1]
	}
	boundaries = boundaries[1:len(boundaries)] // 頭の0値は不要なので破棄

	rand.Seed(time.Now().UnixNano())
	// 境界検知用のスライス
	result := make([]int, len(itemWeights))
	// 結果格納用
	var results []int //抽選結果を入力時のindexで示す

	for i := 0; i < times; i++ {
		draw := rand.Intn(boundaries[len(boundaries)-1]) + 1
		fmt.Println(draw)
		for i, boundary := range boundaries {
			if draw <= boundary {
				result[i]++
				results = append(results, i)
				break
			}
		}
	}
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

	//上の書き方だと読み取れなかった。Queryは複数行、QueryRowが一行らしい
	for row.Next() {
		err := row.Scan(&id, &name, &gtype, &weightRarity[0],
			&weightRarity[1], &weightRarity[2], &weightRarity[3], &weightRarity[4], &start, &end)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(weightRarity)

	//回数だけ抽選を実行する
	var rarityResults []int
	rarityResults = selectWeightRandom(weightRarity, times)
	//e.g. results = {1, 1, 2, 5, 1}

	fmt.Println(rarityResults)

	return results, nil

	for i := 0; i < len(rarityResults); i++ {
		//DBからキャラ抽選の重みとデータを読み込む
		row, err = db.Query("select `characterid`, `weight` from gogame_db.gacha_characters_01_table where `rarity` = ?", rarityResults[id])
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()

		var chara_ids []int
		var weights []int
		j := 0
		for row.Next() {
			if err := row.Scan(&chara_ids[j], &weights[j]); err != nil {
				log.Fatal(err)
			}
			//fmt.Println(id, weight)
		}

		if err := row.Err(); err != nil {
			log.Fatal(err)
		}

		//抽選を実行する

		//抽選結果の名前を問い合わせる

		//抽選結果を格納する
		gachaResult := GachaResult{CharacterId: "", Name: ""}
		results = append(results, gachaResult)
	}

	return results, nil
}
