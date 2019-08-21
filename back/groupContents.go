package back

import (
	"api_test/db"
)

//該当グループの持つ動画を全て取得
func GroupContents(gID uint) (videos []db.Video) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Order("title").Where("group_id = ?", gID).Preload("Channel").Find(&videos)

	return
}
