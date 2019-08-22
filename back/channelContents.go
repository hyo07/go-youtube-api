package back

import (
	"api_test/db"
)

//該当チャンネルの持つ動画を全て取得
func ChannelContents(chID string, viName string, page string) (videos []db.Video) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	offset := Page2Offset(page)
	database.Order("title").
		Where("channel_id = ? AND title LIKE ?", chID, "%"+viName+"%").
		Offset(offset).Limit(10).
		Preload("Channel").Find(&videos)

	return
}
