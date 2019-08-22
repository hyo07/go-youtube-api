package back

import (
	"api_test/db"
)

//該当グループの持つ動画を全て取得
func GroupContents(gID uint, viName string, page string) (videos []db.Video) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	offset := Page2Offset(page)
	database.Order("title").
		Where("group_id = ? AND title LIKE ?", gID, "%"+viName+"%").
		Offset(offset).Limit(10).
		Preload("Channel").Find(&videos)

	return
}
