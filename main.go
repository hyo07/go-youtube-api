package main

import (
	"flag"
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
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

	fmt.Println("検索ワード： " + searchWord + "\n")
	chID := SearchChannelID(searchWord)
	plID := GetSingPlaylist(chID)
	GetPlaylistContnt(plID)

}

func SummarySearchWord(words []string) string {
	var re string
	for _, s := range words {
		re += s + " "
	}
	return re
}

func SearchChannelID(name string) string {
	developerKey := os.Getenv("youtube_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}

	call := service.Search.List("id,snippet").Q(name).MaxResults(10)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	for _, item := range resp.Items {
		if item.Id.Kind == "youtube#channel" {
			return item.Snippet.ChannelId
		}
	}
	panic("Missed get Channel")

}

func GetSingPlaylist(channelID string) string {
	developerKey := os.Getenv("youtube_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}

	call := service.Playlists.List("id,snippet,contentDetails,status").ChannelId(channelID).MaxResults(50)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	itemsConts := 0
	var playlistID string

	for _, playlist := range resp.Items {
		if strings.Contains(playlist.Snippet.Title, "歌") {
			if (!strings.Contains(playlist.Snippet.Title, "配信")) &&
				(!strings.Contains(playlist.Snippet.Title, "枠")) &&
				(!strings.Contains(playlist.Snippet.Title, "放送")) {
				if int(playlist.ContentDetails.ItemCount) > itemsConts {
					playlistID = playlist.Id
				}
			}
		}
	}

	if playlistID == "" {
		panic("Missed get PlaylistID")
	}

	return playlistID

}

func GetPlaylistContnt(playlistID string) {
	developerKey := os.Getenv("youtube_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}

	call := service.PlaylistItems.List("id,snippet,contentDetails,status").PlaylistId(playlistID).MaxResults(50)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	for _, item := range resp.Items {
		fmt.Println(item.Snippet.Title)
		fmt.Println("https://www.youtube.com/watch?v=" + item.ContentDetails.VideoId)
		fmt.Println()
	}

}
