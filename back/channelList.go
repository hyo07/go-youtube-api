package back

import (
	"api_test/db"
)

//DB上の全てのチャンネル取得
//func ChannelList() (channels []db.Channel) {
//	database := db.ConnectDB()
//	defer database.Close()
//	database.LogMode(true)
//
//	database.Order("name").Preload("Group").Find(&channels)
//
//	return
//}

func ChannelSearchList(chName string, page string) (channels []db.Channel) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	offset := Page2Offset(page)
	database.
		Where("name LIKE ?", "%"+chName+"%").
		Order("name").
		Offset(offset).Limit(10).
		Preload("Group").Find(&channels)

	return
}

func ChannelRandomList() (channels []db.Channel) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Order("RANDOM()").Limit(10).Preload("Group").Find(&channels)

	return
}
