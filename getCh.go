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

	db.AddChannel(chContent.ID, chContent.GroupID, chContent.Name, chContent.Thumbnail)
}

//APIクライアント生成
func getClient1() (service *youtube.Service) {
	developerKey := os.Getenv("youtube_key")
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}
	return
}

//チャンネル情報を取得
func getChannelContent1(service *youtube.Service, channelID string, gID uint) (chContent db.Channel) {
	call := service.Channels.List("snippet,contentDetails").Id(channelID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	chContent = db.Channel{
		ID:        channelID,
		GroupID:   gID,
		Name:      resp.Items[0].Snippet.Title,
		Thumbnail: resp.Items[0].Snippet.Thumbnails.Default.Url,
	}

	return
}
