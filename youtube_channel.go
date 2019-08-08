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

	//call := service.Channels.List("id,snippet,contentDetails").Id("UCQ0UDLQCjY0rmuxCDE38FGg").MaxResults(1)
	call := service.Channels.List("id,snippet,contentDetails").Id("UCQ0UDLQCjY0rmuxCDE38FGg").MaxResults(1)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	} else {
		//fmt.Print(resp)
		//fmt.Print(resp.Items[0])
		//fmt.Print(resp.Items[0].Snippet.PublishedAt)

		fmt.Print(resp.Items[0].ContentDetails.RelatedPlaylists.Uploads)

	}
}
