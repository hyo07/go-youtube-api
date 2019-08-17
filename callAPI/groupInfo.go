package callAPI

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
)

func GroupInfo(gID uint) db.Group {
	database, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.Close()
	database.LogMode(true)

	var groups db.Group
	database.Where("id = ?", gID).Preload("Channel").First(&groups)

	return groups
}
