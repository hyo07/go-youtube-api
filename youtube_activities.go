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
	query      = flag.String("query", "夏色まつり", "Search term")
	maxResults = flag.Int64("max-results", 1, "Max YouTube results")
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

	call := service.Search.List("snippet").Q(*query).MaxResults(*maxResults)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	} else {
		//fmt.Print(resp)
		//fmt.Print(resp.Items)
		//fmt.Print(resp.Items[0].Snippet)
		fmt.Print(resp.Items[0].Id.ChannelId)
		//fmt.Print(resp.Items[0].Snippet.Thumbnails.Default.Url)

		//fmt.Print(resp.Items[0].)

		//re := resp.Items[0]
		//fmt.Print(*re)

		//fmt.Print(reflect.TypeOf(resp.Items[0].Snippet))

	}
}
