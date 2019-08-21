package back

import (
	"api_test/db"
)

//DB上の全てのチャンネル取得
func ChannelList() (channels []db.Channel) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Order("name").Preload("Group").Find(&channels)

	return
}

func ChannelSearchList(chName string) (channels []db.Channel) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Order("name").Where("name LIKE ?", "%"+chName+"%").Preload("Group").Find(&channels)

	return
}

func ChannelRandomList() (channels []db.Channel) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.Order("RANDOM()").Preload("Group").Find(&channels)

	return
}
