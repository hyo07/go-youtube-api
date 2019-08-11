package main

/*
県枠ワードから動画を検索して、該当チャンネルから出されている「歌ってみた」動画を抽出
・・・したかったが、「歌」・「歌ってみた」で検索すると、関係ない動画も大量に出てしまい、このやり方で動画を集めるのは厳しいか？
*/

import (
	"flag"
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()

	var searchWord string

	if len(args) == 0 {
		searchWord = "夏色まつり"
	} else {
		searchWord = ssummarySearchWord(args)
	}

	fmt.Println("検索ワード： " + searchWord)
	fmt.Println("----------------------------------")
	fmt.Println()

	chID := ssearchChannelID(searchWord)
	ggetOnluSearchVideo(chID)

}

//コマンドラインで受け取った引数を検索用の文字列に変換
func ssummarySearchWord(words []string) string {
	var re string
	for _, s := range words {
		re += s + " "
	}
	return re
}

//検索用ワードからチャンネルを検索
func ssearchChannelID(name string) string {
	developerKey := os.Getenv("youtube_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}

	call := service.Search.List("snippet").Q(name).MaxResults(1).Type("channel")
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	return resp.Items[0].Snippet.ChannelId
}

func ggetOnluSearchVideo(channelID string) string {
	developerKey := os.Getenv("youtube_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}

	call := service.Search.List("id,snippet").Q("歌ってみた").MaxResults(50).ChannelId(channelID).Type("video")
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	for _, item := range resp.Items {
		//非公開動画は除外・時間で動画のみ抽出
		if (item.Snippet.Title != "Private video") && (ccheckVideoTime(ffilteringSingVideo(item.Id.VideoId))) {
			fmt.Println(item.Snippet.Title)
			fmt.Println("https://www.youtube.com/watch?v=" + item.Id.VideoId)
			fmt.Println("再生時間: " + ffilteringSingVideo(item.Id.VideoId))
			fmt.Println()
		}
	}
	return "0"
}

//各動画の中身を漁る
func ffilteringSingVideo(videoID string) string {
	developerKey := os.Getenv("youtube_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}

	call := service.Videos.List("contentDetails").Id(videoID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	return resp.Items[0].ContentDetails.Duration
}

//歌ってみた「動画」かどうかの判別
func ccheckVideoTime(videoTime string) bool {
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
