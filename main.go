package main

import (
	"api_test/back"
	"api_test/callAPI"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	//TOP
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "top.html", gin.H{})
	})

	//INDEX
	//全動画からランダムで表示
	router.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{"videos": back.RandomVideos()})
	})

	//新しく動画を追加
	//Todo /addにPOSTで飛んじゃってるから、リダイレクトさせて、かつメッセージも渡したい
	router.POST("/add", func(ctx *gin.Context) {
		status, message := callAPI.GetVideo(ctx.PostForm("video-url"))
		if status {
			ctx.HTML(200, "add.html", gin.H{"message": message})
		} else {
			ctx.HTML(200, "add.html", gin.H{"message": message})
		}
	})

	//チャンネル情報と、そのチャンネルの動画を全て表示
	router.GET("/channel/:chID", func(ctx *gin.Context) {
		videos := back.ChannelContents(ctx.Param("chID"))
		channel := back.ChannelInfo(ctx.Param("chID"))
		ctx.HTML(200, "channel.html", gin.H{"videos": videos, "channel": channel})
	})

	//グループ情報と、そのグループに所属するチャンネル一覧。また、そのチャンネルらの持つ動画一覧
	router.GET("/group/:gID", func(ctx *gin.Context) {
		n := ctx.Param("gID")
		gID, err := strconv.Atoi(n)
		if err != nil {
			panic("/group/" + n)
		}
		videos := back.GroupContents(uint(gID))
		group := back.GroupInfo(uint(gID))
		ctx.HTML(200, "group.html", gin.H{"videos": videos, "group": group})
	})

	//チャンネルのリスト
	router.GET("/channels", func(ctx *gin.Context) {
		ctx.HTML(200, "listCh.html", gin.H{"channels": back.ChannelList()})
	})

	//グループのリスト
	router.GET("/groups", func(ctx *gin.Context) {
		ctx.HTML(200, "listGr.html", gin.H{"groups": back.GroupList()})
	})

	router.Run()
}
