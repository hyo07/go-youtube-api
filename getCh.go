package main

/*
   チャンネルIDからチャンネルデータを取り、DBに保存するプログラムの試行
*/

import (
	"api_test/db"
	"flag"
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
)

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

	service := getClient1()
	chContent := getChannelContent1(service, chID, 1)

	fmt.Println(chContent)
	db.InsertCh(chContent.chID, chContent.gID, chContent.name, chContent.thumbnail)
}

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