package callAPI

import "net/url"

func Url2viID(inputURL string) (videoID string) {
	u, err := url.Parse(inputURL)
	if (err != nil) || u.Host != "www.youtube.com" {
		return "URLが正しくありません"
	}
	for k, v := range u.Query() {
		if k == "v" {
			videoID = v[0]
		}
	}
	return
}
