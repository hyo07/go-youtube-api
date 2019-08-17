package back

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
)

//DB上のグループを全て取得
func GroupList() []db.Group {
	database, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.Close()
	database.LogMode(true)

	var groups []db.Group
	database.Preload("Group").Find(&groups)

	return groups
}
