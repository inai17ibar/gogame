package model

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/inai17ibar/gogame/tool"

	_ "github.com/go-sql-driver/mysql"
)

type GachaResult struct {
	CharacterId string `json:"characterID"`
	Name        string `json:"name"`
}

type DrawGachaRequest struct {
	Times int `json:"times"`
}

type DrawGachaResponse struct {
	Results []GachaResult `json:"results"`
}

// dbインスタンス格納用
var db *sql.DB

type LotteryItem struct {
	Id     int
	Weight int
}

func selectWeightRandom(items []LotteryItem, times int) []int {
	// 重みの昇順でチェックしていくので重み一覧をソートする
	sort.SliceStable(items, func(i, j int) bool { return items[i].Weight < items[j].Weight })
	//timesが0でもだめ

	// 抽選結果チェックの基準となる境界値を生成
	boundaries := make([]int, len(items)+1)
	indexs := make([]int, len(items)+1)
	for i := 1; i < len(items)+1; i++ {
		boundaries[i] = boundaries[i-1] + items[i-1].Weight
		indexs[i] = items[i-1].Id
	}
	boundaries = boundaries[1:len(boundaries)] // 頭の0値は不要なので破棄
	indexs = indexs[1:len(indexs)]

	rand.Seed(time.Now().UnixNano())
	// 境界検知用のスライス
	result := make([]int, len(items))
	// 結果格納用
	var results []int //抽選結果をidで示す

	for i := 0; i < times; i++ {
		draw := rand.Intn(boundaries[len(boundaries)-1]) + 1
		//fmt.Println(draw)
		for j, boundary := range boundaries {
			if draw <= boundary {
				result[j]++
				results = append(results, indexs[j])
				break
			}
		}
	}
	// fmt.Printf("境界値: %v\n", boundaries)
	// fmt.Println("重み  想定      結果")
	// fmt.Println("-----------------------------")
	// for i := 0; i < len(items); i++ {
	// 	fmt.Printf("%4d  %f  %f\n",
	// 		items[i].Weight,
	// 		float64(items[i].Weight)/float64(boundaries[len(boundaries)-1]),
	// 		float64(result[i])/float64(times),
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

	items := []LotteryItem{
		{Id: 1, Weight: 0},
		{Id: 2, Weight: 0},
		{Id: 3, Weight: 0},
		{Id: 4, Weight: 0},
		{Id: 5, Weight: 0},
	}

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
		err := row.Scan(&id, &name, &gtype, &items[0].Weight,
			&items[1].Weight, &items[2].Weight, &items[3].Weight, &items[4].Weight, &start, &end)
		if err != nil {
			panic(err)
		}
	}

	//回数だけ抽選を実行する
	var rarityResults []int
	rarityResults = selectWeightRandom(items, times)
	//e.g. results = {1, 1, 2, 5, 1}

	fmt.Println(rarityResults)

	for i := 0; i < len(rarityResults); i++ {
		//DBからキャラ抽選の重みとデータを読み込む
		row, err = db.Query("select `characterid`, `weight` from gogame_db.gacha_characters_01_table where `rarity` = ?", rarityResults[i])
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()

		var tempId int
		var tempWeight int
		var characterItems []LotteryItem
		for row.Next() {
			if err := row.Scan(&tempId, &tempWeight); err != nil {
				log.Fatal(err)
			}
			characterItems = append(characterItems, LotteryItem{Id: tempId, Weight: tempWeight})
		}

		if err := row.Err(); err != nil {
			log.Fatal(err)
		}
		//fmt.Println(characterItems)

		//抽選を実行する
		var result []int
		result = selectWeightRandom(characterItems, 1)

		//fmt.Println(result[0], rarityResults[i])

		//抽選結果の名前を問い合わせる
		character := GachaResult{CharacterId: strconv.Itoa(result[0]), Name: ""}

		row := db.QueryRow("select name from gogame_db.character_table where `id` = ?", result[0])
		err = row.Scan(&character.Name)
		if err := row.Err(); err != nil {
			log.Fatal(err)
		}

		//抽選結果を格納する
		results = append(results, character)
	}

	//fmt.Println(results)

	return results, nil
}
