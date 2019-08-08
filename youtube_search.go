package main

import (
	"flag"
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
)

var (
	query      = flag.String("query", "にじさんじ", "Search term")
	maxResults = flag.Int64("max-results", 5, "Max YouTube results")
)

func main() {
	developerKey := os.Getenv("youtube_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}

	call := service.Search.List("id,snippet").Q(*query).MaxResults(*maxResults)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	for _, item := range resp.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	printIDs("Videos", videos)
	printIDs("Channels", channels)
	printIDs("Playlists", playlists)

}

func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}

