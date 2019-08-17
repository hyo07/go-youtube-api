package back

import (
	"api_test/db"
)

//DB上の全動画の中から、ランダムで一定数を取得
func RandomVideos() []db.Video {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var videos []db.Video
	database.Order("RANDOM()").Limit(5).Preload("Channel").Find(&videos)

	return videos
}

func RandomSearchVideos(viName string) []db.Video {
	database := db.ConnectDB()
	defer database.Close()
	database.LogMode(true)

	var videos []db.Video
	database.Where("title LIKE ?", "%"+viName+"%").Order("RANDOM()").Limit(5).Preload("Channel").Find(&videos)

	return videos
}
