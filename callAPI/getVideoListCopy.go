package callAPI

//
///*
//与えられたチャンネルにある再生リストから歌ってみたのリストを抽出し、保存
//また、そのチャンネルがDBに存在しない場合、追加
//*/
//
//import (
//	"api_test/db"
//	"google.golang.org/api/youtube/v3"
//	"strings"
//)
//
//var reS []map[string]string
//
//func ReadList(chID string) (bool, string) {
//	service := GetClient()
//	plID := getSingPlaylist(service, chID)
//	if len(plID) == 0 {
//		return false, "歌ってみたプレイリストが見つかりませんでした"
//	}
//	for _, pl := range plID {
//		getPlaylistContnt(service, pl)
//	}
//	ret := reS
//	groupID := db.SearchGroup(chID)
//	if !db.CheckExistGroup(groupID) {
//		return false, "与えられたグループが存在しません"
//	}
//
//	for _, video := range ret {
//		switch db.CheckExistVideo(video["channelID"], video["videoID"]) {
//		case 1:
//			db.InsertVideo(video["videoID"], video["channelID"], groupID, video["title"])
//		case 2:
//			chCon3 := GetChannelContent(service, video["channelID"], groupID)
//			db.AddChannel(chCon3.ID, chCon3.GroupID, chCon3.Name, chCon3.Thumbnail)
//			db.InsertVideo(video["videoID"], video["channelID"], groupID, video["title"])
//		default:
//		}
//	}
//	return true, "収集が終わりました"
//}
//
////歌ってみた動画の再生リストを抽出
//func getSingPlaylist(service *youtube.Service, channelID string) (playlistID []string) {
//	var (
//		nextPageToken string
//		plIndex       int64
//		asari         func(ind int64, token string)
//		itemsConts    int64
//	)
//
//	asari = func(ind int64, token string) {
//		call := service.Playlists.List("snippet,contentDetails").ChannelId(channelID).MaxResults(50).PageToken(token)
//		resp, err := call.Do()
//		if err != nil {
//			panic(err)
//		}
//
//		if ind == 0 {
//			plIndex = resp.PageInfo.TotalResults
//			itemsConts = 0
//		}
//		nextPageToken = resp.NextPageToken
//
//		//再生リストの名前で「歌ってみた」動画を判別
//		for _, playlist := range resp.Items {
//			if (strings.Contains(playlist.Snippet.Title, "歌")) ||
//				(strings.Contains(playlist.Snippet.Title, "うた")) ||
//				(strings.Contains(playlist.Snippet.Title, "sing")) ||
//				(strings.Contains(playlist.Snippet.Title, "Music")) {
//				playlistID = append(playlistID, playlist.Id)
//			}
//		}
//		if plIndex > 50 {
//			plIndex -= 50
//			asari(plIndex, nextPageToken)
//		}
//	}
//	asari(0, nextPageToken)
//
//	return
//}
//
////プレイリストの中身を漁る
//func getPlaylistContnt(service *youtube.Service, playlistID string) {
//
//	var (
//		content       map[string]string
//		nextPageToken string
//		plIndex       int64
//		asari         func(ind int64, token string)
//	)
//
//	asari = func(ind int64, token string) {
//		call := service.PlaylistItems.List("snippet,contentDetails").PlaylistId(playlistID).MaxResults(50).PageToken(token)
//		resp, err := call.Do()
//		if err != nil {
//			panic(err)
//		}
//		if ind == 0 {
//			plIndex = resp.PageInfo.TotalResults
//		}
//		nextPageToken = resp.NextPageToken
//
//		for _, item := range resp.Items {
//			//非公開動画は除外・時間で動画のみ抽出
//			if item.Snippet.Title != "Private video" {
//				videoDuration, videoChId := filteringSingVideo(service, item.ContentDetails.VideoId)
//				if CheckVideoTime(videoDuration) {
//					content = map[string]string{
//						"title":     item.Snippet.Title,
//						"videoID":   item.ContentDetails.VideoId,
//						"channelID": videoChId,
//					}
//					reS = append(reS, content)
//				}
//			}
//		}
//		if plIndex > 50 {
//			plIndex -= 50
//			asari(plIndex, nextPageToken)
//		}
//	}
//	asari(0, nextPageToken)
//
//}
//
////各動画の中身を漁る
//func filteringSingVideo(service *youtube.Service, videoID string) (string, string) {
//	call := service.Videos.List("contentDetails").Id(videoID)
//	resp, err := call.Do()
//	if err != nil {
//		panic(err)
//	}
//
//	return resp.Items[0].ContentDetails.Duration, resp.Items[0].Snippet.ChannelId
//}
