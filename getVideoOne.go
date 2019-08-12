package main

/*
投げられたURLの動画を保存する。
また、そのチャンネルがDBに存在しない場合、追加
また、すでにある動画は追加されない
*/

import (
	"api_test/db"
	"flag"
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//動画情報
type videoContent struct {
	viID  string
	chID  string
	gID   uint
	title string
}

//チャンネル情報
type channelContent2 struct {
	chID      string
	gID       uint
	name      string
	thumbnail string
}

func main() {
	flag.Parse()
	args := flag.Args()

	u, err := url.Parse(args[0])
	if err != nil {
		panic("ERRORRRORROR")
	}

	var videoID string
	for k, v := range u.Query() {
		if k == "v" {
			videoID = v[0]
		}
	}

	service := getClient2()
	viContent := getVideoContent2(service, videoID)

	fmt.Println(viContent)

	switch db.CheckExistVideo(viContent.chID, viContent.viID) {
	case 1:
		fmt.Println("チャンネルを確認・動画非重複")
		db.InsertVideo(viContent.viID, viContent.chID, viContent.gID, viContent.title)
	case 2:
		fmt.Println("チャンネルを非確認")
		chCon2 := getChannelContent2(service, viContent.chID, viContent.gID)
		db.AddChannel(chCon2.chID, chCon2.gID, chCon2.name, chCon2.thumbnail)
		db.InsertVideo(viContent.viID, chCon2.chID, chCon2.gID, viContent.title)
	default:
		fmt.Println("False!")
	}
}

//APIクライアント生成
func getClient2() *youtube.Service {
	developerKey := os.Getenv("youtube_key")
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}
	return service
}

//各動画の中身を漁る
func getVideoContent2(service *youtube.Service, videoID string) videoContent {
	call := service.Videos.List("snippet,contentDetails").Id(videoID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	var viContent videoContent

	//非公開動画は除外・時間で動画のみ抽出
	if (resp.Items[0].Snippet.Title != "Private video") && (checkVideoTime2(resp.Items[0].ContentDetails.Duration)) {

		groupID := db.SearchGroup(resp.Items[0].Snippet.ChannelId)

		viContent = videoContent{
			viID:  videoID,
			chID:  resp.Items[0].Snippet.ChannelId,
			gID:   groupID,
			title: resp.Items[0].Snippet.Title,
		}
	}

	return viContent
}

//チャンネル情報を取得
func getChannelContent2(service *youtube.Service, channelID string, gID uint) channelContent2 {
	call := service.Channels.List("snippet,contentDetails").Id(channelID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	chContent2 := channelContent2{
		chID:      channelID,
		gID:       gID,
		name:      resp.Items[0].Snippet.Title,
		thumbnail: resp.Items[0].Snippet.Thumbnails.Default.Url,
	}

	return chContent2
}

//歌ってみた「動画」かどうかの判別
func checkVideoTime2(videoTime string) bool {
	//PT?M?S という形なので、まずPTを消す
	trimTime := strings.TrimLeft(videoTime, "PT")
	//１時間以上のものを除外
	if strings.Contains(trimTime, "H") {
		return false
	}
	//10分以上のものを除外
	if strings.Contains(trimTime, "M") {
		indexM := strings.Index(trimTime, "M")
		minute, _ := strconv.Atoi(trimTime[:indexM])

		if minute >= 10 {
			return false
		}
	}
	
	return true
}
