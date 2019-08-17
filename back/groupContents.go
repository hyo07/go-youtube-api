package back

import (
	"api_test/db"
)

//該当グループの持つ動画を全て取得
func GroupContents(gID uint) []db.Video {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var videos []db.Video
	database.Where("group_id = ?", gID).Preload("Channel").Find(&videos)

	return videos
}
