package back

import (
	"api_test/db"
)

//DB上のグループを全て取得
func GroupList() []db.Group {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var groups []db.Group
	database.Find(&groups)

	return groups
}
