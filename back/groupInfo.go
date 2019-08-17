package back

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
)

//グループ情報取得
func GroupInfo(gID uint) db.Group {
	database, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.Close()
	database.LogMode(true)

	var groups db.Group
	database.Preload("Channel").First(&groups)

	return groups
}
