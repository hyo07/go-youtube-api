package callAPI

/*
   チャンネルIDからチャンネルデータを取り、DBに保存するプログラムの試行
*/

import (
	"api_test/db"
)

func GetChannel(inputURL string, gID uint) (bool, string) {
	chID := Url2chID(inputURL)
	if chID == "1" {
		return false, "youtubeのURLではありません"
	} else if chID == "2" {
		return false, "チャンネルのURLではありません"
	}

	if !db.CheckExistChannel(chID) {
		return false, "既にチャンネルが存在します"
	}

	if !db.CheckExistGroup(gID) {
		return false, "与えられたグループが存在しません"
	}

	service := GetClient()
	chContent := GetChannelContent(service, chID, gID)
	db.AddChannel(chContent.ID, chContent.GroupID, chContent.Name, chContent.Thumbnail)

	return true, "チャンネルを追加しました"
}
