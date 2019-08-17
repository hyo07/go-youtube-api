package back

import (
	"api_test/db"
)

//チャンネル情報取得
func ChannelInfo(chID string) db.Channel {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var channel db.Channel
	database.Where("id = ?", chID).Preload("Group").First(&channel)

	return channel
}
