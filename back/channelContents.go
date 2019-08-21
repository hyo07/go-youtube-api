package back

import (
	"api_test/db"
)

//該当チャンネルの持つ動画を全て取得
func ChannelContents(chID string) (videos []db.Video) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Where("channel_id = ?", chID).Order("title").Preload("Channel").Find(&videos)

	return
}
