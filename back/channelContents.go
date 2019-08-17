package back

import (
	"api_test/db"
)

//該当チャンネルの持つ動画を全て取得
func ChannelContents(chID string) []db.Video {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var videos []db.Video
	database.Where("channel_id = ?", chID).Preload("Channel").Find(&videos)

	return videos
}
