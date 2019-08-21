package callAPI

import (
	"api_test/db"
	"google.golang.org/api/youtube/v3"
)

//チャンネル情報を取得
func GetChannelContent(service *youtube.Service, channelID string, gID uint) (chContent2 db.Channel) {
	call := service.Channels.List("snippet,contentDetails").Id(channelID)
	resp, err := call.Do()
	if err != nil {
		panic(err)
	}

	chContent2 = db.Channel{
		ID:        channelID,
		GroupID:   gID,
		Name:      resp.Items[0].Snippet.Title,
		Thumbnail: resp.Items[0].Snippet.Thumbnails.Default.Url,
	}

	return
}
