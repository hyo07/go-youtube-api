package back

import (
	"api_test/db"
)

//チャンネル情報取得
func ChannelInfo(chID string) (channel db.Channel) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Where("id = ?", chID).Preload("Group").First(&channel)

	return
}
