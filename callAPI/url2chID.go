package callAPI

import (
	"net/url"
	"strings"
)

func Url2chID(inputURL string) string {
	u, err := url.Parse(inputURL)
	if (err != nil) || u.Host != "www.youtube.com" {
		return "1"
	}
	slPa := strings.Split(u.Path, "/")
	if slPa[1] != "channel" {
		return "2"
	}
	var chID string
	chID = slPa[2]
	return chID
}
