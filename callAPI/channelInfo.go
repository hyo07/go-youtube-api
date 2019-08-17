package callAPI

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
)

func ChannelInfo(chID string) db.Channel {
	database, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.Close()
	database.LogMode(true)

	var channel db.Channel
	database.Where("id = ?", chID).Preload("Group").First(&channel)

	return channel
}
