package back

import (
	"api_test/db"
)

//グループ情報取得
func VideoInfo(vID string) db.Video {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var video db.Video
	database.Where("id = ?", vID).Preload("Channel").Preload("Group").First(&video)

	return video
}
