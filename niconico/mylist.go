package niconico

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type MylistVideoThumbnailStyle struct {
	OffsetX int `json:"offset_x"`
	OffsetY int `json:"offset_y"`
	Width   int `json:"width"`
}

type MylistVideo struct {
	ThumbnailURL      string                    `json:"thumbnailurl"`
	Length            string                    `json:"length"`
	LengthSeconds     int                       `json:"length_seconds"`
	Title             string                    `json:"title"`
	ViewCounter       int                       `json:"view_counter"`
	NumRes            int                       `json:"num_res"`
	MylistCounter     int                       `json:"mylist_counter"`
	FirstRetrieve     string                    `json:"first_retrieve"`
	DescriptionShort  string                    `json:"description_short"`
	LastResBody       string                    `json:"last_res_body"`
	ThumbnailStyle    MylistVideoThumbnailStyle `json:"thumbnail_style"`
	IsMiddleThumbnail bool                      `json:"is_middle_thumbnail"`
	ID                string                    `json:"id"`
	CreateTime        int                       `json:"create_time"`
	ThreadUpdateTime  string                    `json:"thread_update_time"`
	MylistComment     string                    `json:"mylist_comment"`
}

type Mylist struct {
	Name                 string        `json:"name"`
	Description          string        `json:"description"`
	UserID               int           `json:"user_id"`
	UserNickname         string        `json:"user_nickname"`
	DefaultSort          int           `json:"default_sort"`
	List                 []MylistVideo `json:"list"`
	IsWatchingThisMylist bool          `json:"is_watching_this_mylist"`
	IsWatchingCountFull  bool          `json:"is_watching_count_full"`
	Status               string        `json:"status"`
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
