package callAPI

/*
DB上のチャンネルのチャンネル情報と、プレイリストを更新する
日〜土の曜日によって1~7の数値が与えられ、その値でDB全体を割ったときの割合部分にあたるチャンネルを更新する。
例：月曜日(=2)、総チャンネル数100
100/7 = 14, 14*2 = 28
よって、先頭から15~28のチャンネルらを更新する
*/

import (
	"api_test/db"
	"github.com/jinzhu/gorm"
	"google.golang.org/api/youtube/v3"
	"time"
)

func UpdateChannel() {
	updateChs, database := weekChannels()
	service := GetClient()
	for i, _ := range updateChs {
		chUpdate(service, updateChs[i], database)
		ReadList(updateChs[i].ID)
	}
}

func weekChannels() ([]db.Channel, *gorm.DB) {
	wdays := [...]int{1, 2, 3, 4, 5, 6, 7}

	t := time.Now()
	weekNum := wdays[t.Weekday()]

	database := db.ConnectDB()
	defer database.Close()

	var channels []db.Channel
	var count int

	database.Select("id").Model(&channels).Count(&count)
	sli := count / 7
	offset := sli * (weekNum - 1)
	if weekNum != 7 {
		database.Offset(offset).Limit(sli).Find(&channels)
	} else {
		database.Offset(offset).Limit(count - offset).Find(&channels)
	}
	return channels, database
}

func chUpdate(service *youtube.Service, channel db.Channel, database *gorm.DB) {
	call := service.Channels.List("snippet,contentDetails").Id(channel.ID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}
	channel.Name = resp.Items[0].Snippet.Title
	channel.Thumbnail = resp.Items[0].Snippet.Thumbnails.Default.Url
	database.Save(&channel)
}
