package main

/*
投げられたURLから取得する
動画からチャンネルを特定し、そのチャンネルがDBにあるか確認する感じを仮予定
*/

import (
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

var reM = make(map[string]string)

func main() {
	flag.Parse()
	args := flag.Args()

	u, err := url.Parse(args[0])
	//u, err := url.Parse("https://www.youtube.com/watch?v=5DfRzFrmHzQ")
	if err != nil {
		panic("ERRORRRORROR")
	}

	var videoID string
	for k, v := range u.Query() {
		if k == "v" {
			videoID = v[0]
		}
	}

	service := getClient()
	chID := getVideoContent(service, videoID)
	getChannelContent(service, chID)

	fmt.Println(reM)
}

func getClient() *youtube.Service {
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
func getVideoContent(service *youtube.Service, videoID string) string {
	call := service.Videos.List("snippet,contentDetails").Id(videoID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	//非公開動画は除外・時間で動画のみ抽出
	if (resp.Items[0].Snippet.Title != "Private video") && (checkVideoTime(resp.Items[0].ContentDetails.Duration)) {
		fmt.Println("Title: " + resp.Items[0].Snippet.Title)
		fmt.Println("再生時間: " + resp.Items[0].ContentDetails.Duration)
		fmt.Println("ビデオID: " + videoID)

		reM["Title"] = resp.Items[0].Snippet.Title
		reM["videoID"] = videoID
		reM["channelID"] = resp.Items[0].Snippet.ChannelId

	}

	return resp.Items[0].Snippet.ChannelId
}

func getChannelContent(service *youtube.Service, channelID string) {
	call := service.Channels.List("snippet,contentDetails").Id(channelID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	fmt.Println("(チャンネル名： " + resp.Items[0].Snippet.Title + ")")
	fmt.Println("チャンネルID： " + channelID)
	fmt.Println("サムネ： " + resp.Items[0].Snippet.Thumbnails.Default.Url)
	fmt.Println()

	//fmt.Println(resp.Items[0].Snippet.Description)
}

//歌ってみた「動画」かどうかの判別
func checkVideoTime(videoTime string) bool {
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
