package main

import (
	"api_test/back"
	"api_test/db"
	"fmt"
)

func main() {
	//fmt.Println(back.Page2Offset("2"))
	database := db.ConnectDB()
	var channels []db.Channel

	offset := back.Page2Offset("2")

	database.Select("name").Offset(offset).Limit(10).Find(&channels)

	for _, ch := range channels {
		fmt.Println(ch.Name)
	}
}
