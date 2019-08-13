package main

/*
   チャンネルIDからチャンネルデータを取り、DBに保存するプログラムの試行
*/

import (
	"api_test/db"
	"flag"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
	"strconv"
)

//チャンネル情報
type channelContent struct {
	chID      string
	gID       uint
	name      string
	thumbnail string
}

func main() {
	flag.Parse()
	args := flag.Args()

	chID := args[0]
	var gID uint
	inputGID, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		panic("第２引数はuintで入力してください")
	}
	gID = uint(inputGID)

	if !db.CheckExistGroup(gID) {
		panic("与えられたグループが存在しません")
	}

	service := getClient1()
	chContent := getChannelContent1(service, chID, gID)

	db.AddChannel(chContent.chID, chContent.gID, chContent.name, chContent.thumbnail)
}

//APIクライアント生成
func getClient1() *youtube.Service {
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

//チャンネル情報を取得
func getChannelContent1(service *youtube.Service, channelID string, gID uint) channelContent {
	call := service.Channels.List("snippet,contentDetails").Id(channelID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	chContent := channelContent{
		chID:      channelID,
		gID:       gID,
		name:      resp.Items[0].Snippet.Thumbnails.Default.Url,
		thumbnail: resp.Items[0].Snippet.Title,
	}

	return chContent
}
