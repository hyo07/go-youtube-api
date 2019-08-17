package back

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
)

//該当グループの持つ動画を全て取得
func GroupContents(gID uint) []db.Video {
	database, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.Close()
	database.LogMode(true)

	var videos []db.Video
	database.Where("group_id = ?", gID).Preload("Channel").Find(&videos)

	return videos
}
