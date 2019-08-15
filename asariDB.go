package main

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.Close()
	database.LogMode(true)

	var videos []db.Video

	//sqliteだと"RANDOM()"、mysqlだと"RAND()" 切り替える
	database.Order("RANDOM()").Limit(5).Find(&videos)
	//database.Order("RANDOM()").Limit(5).Preload("Channel").Find(&videos)

	//database.Raw("SELECT * FROM videos AS vi, " +
	//	"(SELECT id FROM videos ORDER BY RANDOM() LIMIT 0, 5) AS random" +
	//	" WHERE vi.id = random.id LIMIT 0, 5").Scan(&videos)

	//fmt.Println(videos)
	//for _, v := range videos {
	//	fmt.Println(v)
	//}

}
