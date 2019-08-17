package callAPI

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
)

func ChannelContents(chID string) []db.Video {
	database, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.Close()
	database.LogMode(true)

	var videos []db.Video
	database.Where("channel_id = ?", chID).Preload("Channel").Find(&videos)

	return videos
}
