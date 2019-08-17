package callAPI

/*
投げられたURLの動画を保存する。
また、そのチャンネルがDBに存在しない場合、追加
また、すでにある動画は追加されない
*/

import (
	"api_test/db"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"net/url"
)

func GetVideo(inputURL string) (bool, string) {
	u, err := url.Parse(inputURL)
	if (err != nil) || u.Host != "www.youtube.com" {
		return false, "URLが正しくありません"
	}
	var videoID string
	for k, v := range u.Query() {
		if k == "v" {
			videoID = v[0]
		}
	}
	if videoID == "" {
		return false, "URLが正しくありません"
	}

	service := GetClient()
	viContent, err := getVideoContent(service, videoID)
	if err != nil {
		return false, "失敗しました"
	}
	if !db.CheckExistGroup(viContent.GroupID) {
		return false, "与えられたグループが存在しません"
	}

	switch db.CheckExistVideo(viContent.ChannelID, viContent.ID) {
	case 1:
		fmt.Println("チャンネルを確認・動画非重複")
		db.InsertVideo(viContent.ID, viContent.ChannelID, viContent.GroupID, viContent.Title)
		return true, "動画を追加しました"
	case 2:
		fmt.Println("チャンネルを非確認")
		chCon2 := GetChannelContent(service, viContent.ChannelID, viContent.GroupID)
		db.AddChannel(chCon2.ID, chCon2.GroupID, chCon2.Name, chCon2.Thumbnail)
		db.InsertVideo(viContent.ID, chCon2.ID, chCon2.GroupID, viContent.Title)
		return true, "チャンネル・動画を追加しました"
	default:
		return false, "既に追加されています"
	}
}

//各動画の中身を漁る
func getVideoContent(service *youtube.Service, videoID string) (db.Video, error) {
	call := service.Videos.List("snippet,contentDetails").Id(videoID)
	resp, err := call.Do()
	if err != nil {
		return db.Video{}, err
	}

	var viContent db.Video

	//非公開動画は除外・時間で動画のみ抽出
	if (resp.Items[0].Snippet.Title != "Private video") && (CheckVideoTime(resp.Items[0].ContentDetails.Duration)) {

		groupID := db.SearchGroup(resp.Items[0].Snippet.ChannelId)

		viContent = db.Video{
			ID:        videoID,
			ChannelID: resp.Items[0].Snippet.ChannelId,
			GroupID:   groupID,
			Title:     resp.Items[0].Snippet.Title,
		}
	}

	return viContent, nil
}
