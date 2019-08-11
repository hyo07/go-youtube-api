package main

import "fmt"

func main() {
	m1 := map[string]string{
		"Title":     "【誕生日！】愛言葉Ⅲ／夏色まつり cover",
		"channelID": "UCQ0UDLQCjY0rmuxCDE38FGg",
		"videoID":   "jcCGvpvxqVQ",
	}

	m2 := map[string]string{
		"Title":     "カタオモイ(Aimer)  /  夏色まつり cover",
		"channelID": "UCQ0UDLQCjY0rmuxCDE38FGg",
		"videoID":   "TT6gTs2B1Uw",
	}

	m3 := map[string]string{
		"Title":     "【歌ってみた】スイートマジック 【夏色まつり×紫咲シオン】",
		"channelID": "UCQ0UDLQCjY0rmuxCDE38FGg",
		"videoID":   "z2pWWKSPofA",
	}

	m4 := map[string]string{
		"Title":     "【歌ってみた】お願いダーリン ／夏色まつり",
		"channelID": "UCQ0UDLQCjY0rmuxCDE38FGg",
		"videoID":   "kMPx3AHvwtA",
	}

	m5 := map[string]string{
		"Title":     "【歌ってみた】好き！雪！本気マジック ／夏色まつり",
		"channelID": "UCQ0UDLQCjY0rmuxCDE38FGg",
		"videoID":   "nZ-g1pPNfts",
	}

	m6 := map[string]string{
		"Title":     "【歌ってみた】未練レコード／夏色まつり【ときのそら × 40mP】",
		"channelID": "UCQ0UDLQCjY0rmuxCDE38FGg",
		"videoID":   "1mOL7Hi46h8",
	}

	m7 := map[string]string{
		"Title":     "【HoneyWorks】ファンサ／夏色まつり cover",
		"channelID": "UCQ0UDLQCjY0rmuxCDE38FGg",
		"videoID":   "xpUQ-5dSJTk",
	}

	contents := []map[string]string{m1, m2, m3, m4, m5, m6, m7}

	fmt.Println(contents)

}
