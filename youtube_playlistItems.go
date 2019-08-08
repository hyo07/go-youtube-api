package main

import (
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
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

	call := service.PlaylistItems.List("id,snippet,contentDetails,status").PlaylistId("PL6sZ3uYmeG1vGZVLNlWdJkm3W0R86y_Qa").MaxResults(50)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	//fmt.Print(resp.Items[0].Snippet.Title)

	//fmt.Println(reflect.TypeOf(resp.Items))

	for _, item := range resp.Items {
		fmt.Println(item.Snippet.Title)
		fmt.Println("https://www.youtube.com/watch?v=" + item.ContentDetails.VideoId)
		fmt.Println()
	}

}
