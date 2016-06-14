package niconico

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type MylistVideoThumbnailStyle struct {
	OffsetX int
	OffsetY int
	Width   int
}

type MylistVideo struct {
	ThumbnailURL      string
	Length            string
	LengthSeconds     int
	Title             string
	ViewCounter       int
	NumRes            int
	MylistCounter     int
	FirstRetrieve     string
	DescriptionShort  string
	LastResBody       string
	ThumbnailStyle    MylistVideoThumbnailStyle
	IsMiddleThumbnail bool
	ID                string
	CreateTime        int
	ThreadUpdateTime  string
	MylistComment     string
}

type Mylist struct {
	Name                 string
	Description          string
	UserID               int
	UserNickname         string
	DefaultSort          int
	List                 []MylistVideo
	IsWatchingThisMylist bool
	IsWatchingCountFull  bool
	Status               string
}

func IsMylist(query string) bool {
	if m, _ := regexp.MatchString("mylist/", query); !m {
		return false
	}
	return true
}

func ToMylistID(query string) string {
	re, _ := regexp.Compile("\\d+")
	one := re.Find([]byte(query))

	return string(one)
}

func GetMylist(mylistID string, sessionKey string) (Mylist, error) {
	target := "http://riapi.nicovideo.jp/api/watch/mylistvideo?id="
	target += mylistID
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Add("Cookie", sessionKey)
	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()

	byteArray, _ := ioutil.ReadAll(res.Body)

	var m Mylist
	err := json.Unmarshal(byteArray, &m)

	return m, err
}
