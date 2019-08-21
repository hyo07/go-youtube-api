package back

import (
	"api_test/db"
)

//DB上のグループを全て取得
func GroupList() (groups []db.Group) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Find(&groups)

	return
}
