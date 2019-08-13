package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

//チャンネル情報
type Channel struct {
	ID           string `gorm:"primary_key"`
	GroupID      uint
	Group        Group
	Name         string  `gorm:"unique;not null"`
	Thumbnail    string  `gorm:"unique;not null"`
	Descripition string  `gorm:"type:text"`
	Video        []Video `gorm:"foreignkey:ChannelID"`
}

//該当Vtuberのグループ情報
type Group struct {
	ID      uint      `gorm:"primary_key;AUTO_INCREMENT"`
	Name    string    `gorm:"unique;not null"`
	Channel []Channel `gorm:"foreignkey:GroupID"`
	Video   []Video   `gorm:"foreignkey:GroupID"`
}

//動画情報
type Video struct {
	ID           string `gorm:"primary_key"`
	ChannelID    string
	Channel      Channel
	GroupID      uint
	Group        Group
	Title        string `gorm:"unique;not null"`
	Descripition string `gorm:"type:text"`
}

func main() {
	db, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&Channel{}, &Group{}, &Video{})

	//channel := Channel{
	//	ID:           "UCQ0UDLQCjY0rmuxCDE38FGg",
	//	GroupID:      1,
	//	Name:         "Matsuri Channel 夏色まつり",
	//	Thumbnail:    "https://yt3.ggpht.com/a/AGF-l7_MTJEH9Kn-cVznPJPBt4v0BOkmd5btoSdz6Q=s88-c-k-c0xffffffff-no-rj-mo",
	//	Descripition: "チャンネル説明",
	//}
	//db.NewRecord(&channel)
	//db.Create(&channel)
	//
	//group := Group{
	//	ID:   2,
	//	Name: "ホロライブ",
	//}
	//db.NewRecord(&group)
	//db.Create(&group)
	//
	//video := Video{
	//	ID:           "jcCGvpvxqVQ",
	//	ChannelID:    "UCQ0UDLQCjY0rmuxCDE38FGg",
	//	GroupID:      1,
	//	Title:        "【誕生日！】愛言葉Ⅲ／夏色まつり cover",
	//	Descripition: "動画説明",
	//}
	//db.NewRecord(&video)
	//db.Create(&video)

	//var videos Video
	//db.Find(&videos).Related(&videos.Channel).Related(&videos.Group)
	//fmt.Println(videos)

	//fmt.Println(checkExistVideo("UCQ0UDLQCjY0rmuxCDE38FGgあ", "hoge"))
	//fmt.Println(checkExistVideo("UCQ0UDLQCjY0rmuxCDE38FGg", "jcCGvpvxqVQ"))
	//SearchGroup("UCQ0UDLQCjY0rmuxCDE38FGg")

	//ChangeChGroup("UCXTpFs_3PqI41qX2d9tL2Rw", 2)
}

//新しくチャンネルをDBに追加
func AddChannel(chID string, gID uint, name string, thumbnail string) {
	db, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//db.LogMode(true)

	var channel Channel
	db.Where("id = ?", chID).Find(&channel)
	//重複チェック
	if channel.ID != "" {
		panic("既に存在します(channel)")
	} else {
		channel = Channel{
			ID:        chID,
			GroupID:   gID,
			Name:      name,
			Thumbnail: thumbnail,
		}
		db.NewRecord(&channel)
		db.Create(&channel)
		fmt.Println("チャンネルを追加しました")
	}
}

//渡された動画が既にDBにないか確認
func CheckExistVideo(chID string, viID string) int {
	db, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//db.LogMode(true)

	var videos Video
	var channels Channel

	db.Where("id = ?", chID).Find(&channels)
	//チャンネルすら無いのか・チャンネルはあるが動画はないのか・どっちも既にあるのか
	if channels.ID == "" {
		fmt.Println("存在しません(channel)")
		return 2
	}

	db.Where("id = ?", viID).Find(&videos)
	if videos.ID != "" {
		fmt.Println("すでに存在します(video)")
		return 0
	}
	fmt.Println("動画を追加できます")
	return 1
}

//動画をDBに保存
func InsertVideo(viID string, chID string, gID uint, title string) {
	db, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.LogMode(true)

	var videos Video

	videos = Video{
		ID:        viID,
		ChannelID: chID,
		GroupID:   gID,
		Title:     title,
	}
	db.NewRecord(&videos)
	db.Create(&videos)
	fmt.Println("動画を追加しました")
}

//該当チャンネルの所属グループを確認
func SearchGroup(chID string) uint {
	db, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//db.LogMode(true)

	var channel Channel
	db.Where("id = ?", chID).Find(&channel)
	//fmt.Println(channel.GroupID)

	//所属グループが割り振られてい無い場合、1（その他）を入れる
	if channel.GroupID == 0 {
		return 1
	} else {
		return channel.GroupID
	}
}

//該当チャンネル、および投稿動画のグループを変更
func ChangeChGroup(chID string, newGroupID uint) {
	db, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//db.LogMode(true)

	var group Group
	db.Where("id = ?", newGroupID).Find(&group)

	if group.ID != 0 {
		var channel Channel
		var videos []Video

		//チャンネルのGroupIDを変更
		db.Where("id = ?", chID).Take(&channel)
		channel.GroupID = newGroupID
		db.Save(&channel)

		//該当チャンネルの持つ投稿動画を全て変更
		db.Find(&videos).Where("channel_id = ?", chID)
		for i, _ := range videos {
			videos[i].GroupID = newGroupID
			db.Save(&videos[i])
		}
		fmt.Println("グループを変更しました")
	} else {
		fmt.Println("存在しないグループです")
	}

}

//グループの存在を確認
func CheckExistGroup(gID uint) bool {
	db, err := gorm.Open("sqlite3", "./db/test.sqlite3")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//db.LogMode(true)

	var group Group

	db.Where("id = ?", gID).Find(&group)
	if group.ID == 0 {
		return false
	} else {
		return true
	}
}
