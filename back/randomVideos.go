package back

import (
	"api_test/db"
)

//DB上の全動画の中から、ランダムで一定数を取得
//func RandomVideos() (videos []db.Video) {
//	database := db.ConnectDB()
//	defer database.Close()
//	database.LogMode(true)
//
//	database.Order("RANDOM()").Limit(5).Preload("Channel").Find(&videos)
//
//	return
//}

func RandomVideos(viName string) (videos []db.Video) {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	database.
		Where("title LIKE ?", "%"+viName+"%").
		Order("RANDOM()").
		Limit(10).
		Preload("Channel").
		Find(&videos)

	return
}
