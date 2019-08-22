package main

import (
	"api_test/back"
	"api_test/callAPI"
	"api_test/db"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"strconv"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/assets", "./assets")

	//TOP
	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/index"
		router.HandleContext(ctx)
	})

	//INDEX
	//全動画からランダムで表示
	router.GET("/index", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
		router.SetHTMLTemplate(html)

		ctx.HTML(200, "base.html", gin.H{"videos": back.RandomVideos(ctx.Query("viName"))})
	})

	//新しく情報を追加
	router.GET("/add", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/add.html"))
		router.SetHTMLTemplate(html)

		ctx.HTML(200, "base.html", gin.H{"groups": back.GroupList("", "")})
	})

	//ビデオ追加
	router.POST("/add-vi", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/add.html"))
		router.SetHTMLTemplate(html)

		status, message := callAPI.GetVideo(ctx.PostForm("video-url"))
		if status {
			viID := callAPI.Url2viID(ctx.PostForm("video-url"))
			video := back.VideoInfo(viID)
			ctx.HTML(200, "base.html", gin.H{"message": message, "video": video})
		} else {
			//TODO リダイレクトとしては不適切？ /addにmessageだけ送りたいがどうしようか？
			html := template.Must(template.ParseFiles("templates/base.html", "templates/add.html"))
			router.SetHTMLTemplate(html)
			ctx.HTML(200, "base.html", gin.H{"groups": back.GroupList("", ""), "message": message})
		}
	})

	//チャンネル追加
	router.POST("/add-ch", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/added.html"))
		router.SetHTMLTemplate(html)

		inputURL := ctx.PostForm("channel-url")
		gID, err := strconv.Atoi(ctx.PostForm("select-group"))
		if err != nil {
			ctx.HTML(200, "base.html", gin.H{"message": "グループ選択にエラーが発生しました"})
		}
		status, message := callAPI.GetChannel(inputURL, uint(gID))
		if status {
			chID := callAPI.Url2chID(inputURL)
			channel := back.ChannelInfo(chID)
			callAPI.ReadList(chID)
			videos := back.ChannelContents(chID, "", "1")
			ctx.HTML(200, "base.html", gin.H{"message": message, "channel": channel, "videos": videos})
		} else {
			//TODO リダイレクトとしては不適切？ /addにmessageだけ送りたいがどうしようか？
			html := template.Must(template.ParseFiles("templates/base.html", "templates/add.html"))
			router.SetHTMLTemplate(html)
			ctx.HTML(200, "base.html", gin.H{"groups": back.GroupList("", ""), "message": message})
		}
	})

	//ビデオ詳細
	router.GET("/video/:vID", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/video.html"))
		router.SetHTMLTemplate(html)

		video := back.VideoInfo(ctx.Param("vID"))
		ctx.HTML(200, "base.html", gin.H{"video": video})
	})

	//チャンネル情報と、そのチャンネルの動画を全て表示
	//TODO ページ上限取得
	router.GET("/channel/:chID", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/channel.html"))
		router.SetHTMLTemplate(html)
		videos := back.ChannelContents(ctx.Param("chID"), ctx.Query("viName"), ctx.Query("page"))
		channel := back.ChannelInfo(ctx.Param("chID"))
		ctx.HTML(200, "base.html", gin.H{"videos": videos, "channel": channel})
	})

	//グループ情報と、そのグループに所属するチャンネル一覧。また、そのチャンネルらの持つ動画一覧
	//TODO ページ上限取得
	router.GET("/group/:gID", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/group.html"))
		router.SetHTMLTemplate(html)

		n := ctx.Param("gID")
		gID, err := strconv.Atoi(n)
		if err != nil {
			gID = 1
		}
		display := ctx.Query("display")
		var group db.Group
		if display == "channel" {
			group = back.GroupInfo(uint(gID))
			channels := back.GroupHasCh(uint(gID), ctx.Query("chName"), ctx.Query("page"))
			ctx.HTML(200, "base.html", gin.H{"group": group, "channels": channels})
		} else {
			group = back.GroupInfo(uint(gID))
			videos := back.GroupContents(uint(gID), ctx.Query("viName"), ctx.Query("page"))
			ctx.HTML(200, "base.html", gin.H{"videos": videos, "group": group})
		}
	})

	//チャンネルのリスト
	//TODO ページ上限取得
	router.GET("/channels", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/listCh.html"))
		router.SetHTMLTemplate(html)

		chName := ctx.Query("chName")
		if chName == "random" {
			ctx.HTML(200, "base.html", gin.H{"channels": back.ChannelRandomList()})
		} else {
			ctx.HTML(200, "base.html", gin.H{"channels": back.ChannelSearchList(chName, ctx.Query("page"))})
		}
	})

	//グループのリスト
	//TODO ページ上限取得
	router.GET("/groups", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/base.html", "templates/listGr.html"))
		router.SetHTMLTemplate(html)

		gName := ctx.Query("gName")
		if gName == "random" {
			ctx.HTML(200, "base.html", gin.H{"groups": back.GroupRandomList()})
		} else {
			ctx.HTML(200, "base.html", gin.H{"groups": back.GroupList(gName, ctx.Query("page"))})
		}
	})

	router.Run()
}
