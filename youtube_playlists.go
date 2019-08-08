package main

import (
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
	"strings"
)

//var (
//	query      = flag.String("query", "夏色まつり", "Search term")
//	maxResults = flag.Int64("max-results", 1, "Max YouTube results")
//)

func main() {
	developerKey := os.Getenv("youtube_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Errr creating new Youtube client: %v", err)
	}

	call := service.Playlists.List("id,snippet,status").ChannelId("UCQ0UDLQCjY0rmuxCDE38FGg").MaxResults(50)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	//fmt.Print(resp.Items[0].Snippet.Title)

	//fmt.Println(reflect.TypeOf(resp.Items))

	for _, playlist := range resp.Items {
		if strings.Contains(playlist.Snippet.Title, "歌") {
			fmt.Println(">>>>>" + playlist.Snippet.Title + "<<<<<")
			fmt.Println(">>>>>" + playlist.Id + "<<<<<")
		} else {
			fmt.Println(playlist.Snippet.Title)
		}
	}

}
