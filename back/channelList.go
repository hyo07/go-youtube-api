package back

import (
	"api_test/db"
)

//DB上の全てのチャンネル取得
func ChannelList() []db.Channel {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var channels []db.Channel
	database.Order("name").Preload("Group").Find(&channels)

	return channels
}

func ChannelSearchList(chName string) []db.Channel {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var channels []db.Channel
	database.Order("name").Where("name LIKE ?", "%"+chName+"%").Preload("Group").Find(&channels)

	return channels
}
