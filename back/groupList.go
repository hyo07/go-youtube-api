package back

import (
	"api_test/db"
)

//DB上のグループを全て取得
func GroupList(gName string, page string) (groups []db.Group) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	offset := Page2Offset(page)
	database.
		Where("name LIKE ?", "%"+gName+"%").
		Order("id").
		Offset(offset).Limit(10).
		Find(&groups)

	return
}

func GroupRandomList() (groups []db.Group) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Order("RANDOM()").Limit(10).Preload("Group").Find(&groups)

	return
}
