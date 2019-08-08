package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Data struct {
	Url            string `json:"url"`
	Title          string `json:"title"`
	LikesCount     int    `json:"likes_count"`
	ReactionsCount int    `json:"reactions_count"`
	PageViewsCount int    `json:"page_views_count"`
}

func fetchQiitaData(accessToken string, targetDate time.Time) []Data {
	baseUrl := "https://qiita.com/api/v2/"
	action := "items"

	baseParam := "?page=1&per_page=30"

	year, month, day := targetDate.Date()
	targetDay := dateNum2String(year, int(month), day)

	year, month, day = targetDate.AddDate(0, 0, 1).Date()
	nextDay := dateNum2String(year, int(month), day)

	varParam := "&query=stocks:>30+created:>=" + targetDay + "+created:<" + nextDay

	endpointURL, err := url.Parse(baseUrl + action + baseParam + varParam)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(Data{})
	if err != nil {
		panic(err)
	}

	var resp = &http.Response{}
	if len(accessToken) > 0 {
		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL,
			Method: "GET",
			Header: http.Header{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer " + accessToken},
			},
		})
	} else {
		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL,
			Method: "GET",
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		})
	}

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var datas []Data

	if err := json.Unmarshal(b, &datas); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil
	}
	return datas
}

func outputQiitaData(datas []Data) {
	for _, val := range datas {
		fmt.Println(val.LikesCount, val.Title, val.Url)
	}
	fmt.Println()
}

func dateNum2String(year int, month int, day int) string {
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func main() {
	//qiitaToken := os.Getenv("QIIT_TOKEN")

	var baseDate time.Time

	fmt.Println("いいね数       タイトル        URL")

	var err error
	if len(os.Args) >= 2 {
		arg := os.Args[1]
		layout := "2006-01-02"
		baseDate, err = time.Parse(layout, arg)
		if err != nil {
			panic(err)
		}
	} else {
		baseDate = time.Now()
	}
	startDate := baseDate.AddDate(-1, 0, -6)
	num0fDates := 7

	for i := 0; i < num0fDates; i++ {
		targetDate := startDate.AddDate(0, 0, i)
		fmt.Printf("%d-%d-%d\n", targetDate.Year(), int(targetDate.Month()), targetDate.Day())

		datas := fetchQiitaData(qiitaToken, targetDate)

		outputQiitaData(datas)
	}
}
