package main

/*
与えられたチャンネルにある再生リストから歌ってみたのリストを抽出し、保存
また、そのチャンネルがDBに存在しない場合、追加
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
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()

	//var searchWord string
	//
	//if len(args) == 0 {
	//	searchWord = "夏色まつり"
	//} else {
	//	searchWord = summarySearchWord3(args)
	//}
	//
	//fmt.Println("検索ワード： " + searchWord)
	//fmt.Println("----------------------------------")
	//fmt.Println()

	service := getClient3()

	//chID := searchChannelID3(service, searchWord)

	chID := args[0]
	plID := getSingPlaylist3(service, chID)
	ret := getPlaylistContnt3(service, plID)

	fmt.Println(ret)

	groupID := db.SearchGroup(chID)

	if !db.CheckExistGroup(groupID) {
		panic("与えられたグループが存在しません")
	}

	for _, video := range ret {
		fmt.Println("------------------------------------------")
		switch db.CheckExistVideo(video["channelID"], video["videoID"]) {
		case 1:
			fmt.Println("チャンネルを確認・動画非重複")
			db.InsertVideo(video["videoID"], video["channelID"], groupID, video["title"])
		case 2:
			fmt.Println("チャンネルを非確認")
			chCon3 := getChannelContent3(service, video["channelID"], groupID)
			db.AddChannel(chCon3.ID, chCon3.GroupID, chCon3.Name, chCon3.Thumbnail)
			db.InsertVideo(video["videoID"], video["channelID"], groupID, video["title"])
		default:
			fmt.Println("動画追加済み")
		}
		fmt.Println("------------------------------------------")
	}
}

//APIクライアント生成
func getClient3() *youtube.Service {
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
func summarySearchWord3(words []string) string {
	var re string
	for _, s := range words {
		re += s + " "
	}
	return re
}

//検索用ワードからチャンネルを検索
func searchChannelID3(service *youtube.Service, name string) string {
	call := service.Search.List("snippet").Q(name).MaxResults(1).Type("channel")
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}
	return resp.Items[0].Snippet.ChannelId
}

//歌ってみた動画の再生リストを抽出
func getSingPlaylist3(service *youtube.Service, channelID string) string {
	var (
		nextPageToken string
		plIndex       int64
		asari         func(ind int64, token string)
		playlistID    string
		itemsConts    int64
	)

	asari = func(ind int64, token string) {
		call := service.Playlists.List("snippet,contentDetails").ChannelId(channelID).MaxResults(50).PageToken(token)
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
func getPlaylistContnt3(service *youtube.Service, playlistID string) []map[string]string {

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
			if (item.Snippet.Title != "Private video") && (checkVideoTime3(filteringSingVideo3(service, item.ContentDetails.VideoId))) {
				content = map[string]string{
					"title":     item.Snippet.Title,
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
func filteringSingVideo3(service *youtube.Service, videoID string) string {
	call := service.Videos.List("contentDetails").Id(videoID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	return resp.Items[0].ContentDetails.Duration
}

//「歌ってみた動画」かどうかの判別
func checkVideoTime3(videoTime string) bool {
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

//チャンネル情報を取得
func getChannelContent3(service *youtube.Service, channelID string, gID uint) db.Channel {
	call := service.Channels.List("snippet,contentDetails").Id(channelID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	chContent3 := db.Channel{
		ID:        channelID,
		GroupID:   gID,
		Name:      resp.Items[0].Snippet.Title,
		Thumbnail: resp.Items[0].Snippet.Thumbnails.Default.Url,
	}

	return chContent3
}
