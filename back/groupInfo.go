package back

import (
	"api_test/db"
)

//グループ情報取得
func GroupInfo(gID uint) db.Group {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var group db.Group
	database.Where("id = ?", gID).Preload("Channel").First(&group)

	return group
}
