package callAPI

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
)

func ChannelList() []db.Channel {
	database, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.Close()
	database.LogMode(true)

	var channels []db.Channel
	database.Preload("Group").Find(&channels)

	return channels
}
