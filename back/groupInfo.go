package back

import (
	"api_test/db"
)

//グループ情報取得
func GroupInfo(gID uint) (group db.Group) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Where("id = ?", gID).Preload("Channel").First(&group)

	return
}

func GroupInfoNoCh(gID uint) (group db.Group) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Where("id = ?", gID).First(&group)

	return
}
