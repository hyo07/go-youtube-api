package main

/*
Vtuberの名前から、その人の再生リストを見にいき、「歌ってみた」動画の再生リストの中身を取ってくる
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
		searchWord = SummarySearchWord(args)
	}

	fmt.Println("検索ワード： " + searchWord)
	fmt.Println("----------------------------------")
	fmt.Println()

	service := GetClient()
	chID := SearchChannelID(service, searchWord)
	plID := GetSingPlaylist(service, chID)
	GetPlaylistContnt(service, plID)
	//contents := GetPlaylistContnt(service, plID)
	//fmt.Println(contents)
}

func GetClient() *youtube.Service {
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

//コマンドラインで受け取った引数を検索用の文字列に変換
func SummarySearchWord(words []string) string {
	var re string
	for _, s := range words {
		re += s + " "
	}
	return re
}

//検索用ワードからチャンネルを検索
func SearchChannelID(service *youtube.Service, name string) string {
	call := service.Search.List("snippet").Q(name).MaxResults(1).Type("channel")
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}
	return resp.Items[0].Snippet.ChannelId
}

//歌ってみた動画の再生リストを抽出
func GetSingPlaylist(service *youtube.Service, channelID string) string {
	var (
		nextPageToken string
		plIndex       int64
		asari         func(ind int64, token string)
		playlistID    string
		itemsConts    int64
	)

	asari = func(ind int64, token string) {
		call := service.Playlists.List("snippet,contentDetails").ChannelId(channelID).MaxResults(50).PageToken(token)
		fmt.Println("Do!!!!!")
		resp, err := call.Do()
		if err != nil {
			panic(err)
		}

		if ind == 0 {
			plIndex = resp.PageInfo.TotalResults
			itemsConts = 0
		}
		nextPageToken = resp.NextPageToken

		//再生リストの名前で「歌ってみた」動画を判別
		for _, playlist := range resp.Items {
			if (strings.Contains(playlist.Snippet.Title, "歌")) || (strings.Contains(playlist.Snippet.Title, "うた")) {
				if (!strings.Contains(playlist.Snippet.Title, "配信")) &&
					(!strings.Contains(playlist.Snippet.Title, "枠")) &&
					(!strings.Contains(playlist.Snippet.Title, "わく")) &&
					(!strings.Contains(playlist.Snippet.Title, "雑談")) &&
					(!strings.Contains(playlist.Snippet.Title, "放送")) {
					if playlist.ContentDetails.ItemCount > itemsConts {
						playlistID = playlist.Id
					}
				}
			}
		}
		if plIndex > 50 {
			plIndex -= 50
			asari(plIndex, nextPageToken)
		}
	}

	asari(0, nextPageToken)

	if playlistID == "" {
		panic("Missed get PlaylistID")
	}

	return playlistID
}

//プレイリストの中身を漁る
func GetPlaylistContnt(service *youtube.Service, playlistID string) []map[string]string {

	var (
		reS           []map[string]string
		content       map[string]string
		nextPageToken string
		plIndex       int64
		asari         func(ind int64, token string)
	)

	asari = func(ind int64, token string) {
		call := service.PlaylistItems.List("snippet,contentDetails").PlaylistId(playlistID).MaxResults(50).PageToken(token)
		resp, err := call.Do()
		if err != nil {
			panic(err)
		}
		if ind == 0 {
			plIndex = resp.PageInfo.TotalResults
		}
		nextPageToken = resp.NextPageToken

		for _, item := range resp.Items {
			//非公開動画は除外・時間で動画のみ抽出
			if (item.Snippet.Title != "Private video") && (CheckVideoTime(FilteringSingVideo(service, item.ContentDetails.VideoId))) {
				//if item.Snippet.Title != "Private video" {

				fmt.Println(item.Snippet.Title)
				fmt.Println("https://www.youtube.com/watch?v=" + item.ContentDetails.VideoId)
				fmt.Println("再生時間: " + FilteringSingVideo(service, item.ContentDetails.VideoId))
				//fmt.Println("videoID: " + item.ContentDetails.VideoId)
				//fmt.Println("channelID: "+ item.Snippet.ChannelId)
				//fmt.Println("サムネURL: " + item.Snippet.Thumbnails.Default.Url)
				//fmt.Println("インデックス: " + strconv.Itoa(item.Snippet.Position))
				fmt.Println()

				content = map[string]string{
					"Title":     item.Snippet.Title,
					"videoID":   item.ContentDetails.VideoId,
					"channelID": item.Snippet.ChannelId,
				}
				reS = append(reS, content)
			}
		}
		if plIndex > 50 {
			plIndex -= 50
			asari(plIndex, nextPageToken)
		}
	}
	asari(0, nextPageToken)
	fmt.Println("---------------------\n" + "合計取得数： " + strconv.Itoa(len(reS)) + "\n---------------------")

	return reS
}

//各動画の中身を漁る
func FilteringSingVideo(service *youtube.Service, videoID string) string {
	call := service.Videos.List("contentDetails").Id(videoID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	return resp.Items[0].ContentDetails.Duration
}

//歌ってみた「動画」かどうかの判別
func CheckVideoTime(videoTime string) bool {
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
