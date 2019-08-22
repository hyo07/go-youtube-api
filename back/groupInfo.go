package back

import (
	"api_test/db"
)

//グループ情報取得
func GroupInfo(gID uint) (group db.Group) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Where("id = ?", gID).First(&group)

	return
}

func GroupHasCh(gID uint, chName string, page string) (channels []db.Channel) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	offset := Page2Offset(page)
	database.Where("group_id = ? AND name LIKE ?", gID, "%"+chName+"%").
		Order("name").
		Offset(offset).Limit(10).
		Preload("Group").Find(&channels)

	return
}
