package back

import (
	"api_test/db"
)

//ビデオ情報取得
func VideoInfo(vID string) (video db.Video) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Where("id = ?", vID).Preload("Channel").Preload("Group").First(&video)

	return
}
