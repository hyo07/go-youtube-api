package callAPI

/*
   チャンネルIDからチャンネルデータを取り、DBに保存するプログラムの試行
*/

import (
	"api_test/db"
	"net/url"
	"strings"
)

func GetChannel(inputURL string, gID uint) (bool, string) {
	u, err := url.Parse(inputURL)
	if (err != nil) || u.Host != "www.youtube.com" {
		return false, "youtubeのURLではありません"
	}
	slPa := strings.Split(u.Path, "/")
	if slPa[1] != "channel" {
		return false, "チャンネルのURLではありません"
	}
	var chID string
	chID = slPa[2]

	if !db.CheckExistGroup(gID) {
		return false, "与えられたグループが存在しません"
	}
	service := GetClient()
	chContent := GetChannelContent(service, chID, gID)
	db.AddChannel(chContent.ID, chContent.GroupID, chContent.Name, chContent.Thumbnail)

	return true, "チャンネルを追加しました"
}
