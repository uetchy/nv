package niconico

import (
	"encoding/xml"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Thumbinfo struct {
	VideoID       string `xml:"thumb>video_id"`
	Title         string `xml:"thumb>title"`
	Description   string `xml:"thumb>description"`
	ThumbnailURL  string `xml:"thumb>thumbnail_url"`
	FirstRetrieve string `xml:"thumb>first_retrieve"`
	Length        string `xml:"thumb>length"`
	MovieType     string `xml:"thumb>movie_type"`
	SizeHigh      int    `xml:"thumb>size_high"`
	SizeLow       int    `xml:"thumb>size_low"`
	ViewCounter   int    `xml:"thumb>view_counter"`
	CommentNum    int    `xml:"thumb>comment_num"`
	MylistCounter int    `xml:"thumb>mylist_counter"`
	LastResBody   string `xml:"thumb>last_res_body"`
	WatchURL      string `xml:"thumb>watch_url"`
	ThumbType     string `xml:"thumb>thumb_type"`
	Embeddable    string `xml:"thumb>embeddable"`
	NoLivePlay    int    `xml:"thumb>no_live_play"`
	Tags          []Tag  `xml:"thumb>tags>tag"`
	UserID        int    `xml:"thumb>user_id"`
	UserNickname  string `xml:"thumb>user_nickname"`
	UserIconURL   string `xml:"thumb>user_icon_url"`
}

type Tag struct {
	Title    string `xml:",chardata"`
	Category string `xml:"category,attr"`
}

func ToVideoID(query string) string {
	re, _ := regexp.Compile("[a-z]{2}?\\d+")
	one := re.Find([]byte(query))

	return string(one)
}

func GetThumbInfo(videoID string) (thumb Thumbinfo, err error) {
	target := "http://ext.nicovideo.jp/api/getthumbinfo/" + videoID
	res, err := http.Get(target)
	if err != nil {
		return thumb, err
	}
	defer res.Body.Close()

	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return thumb, err
	}
	if err := xml.Unmarshal(byteArray, &thumb); err != nil {
		return thumb, err
	}

	return thumb, nil
}

func GetHistory(videoID string, sessionKey string) (nicoHistory string, err error) {
	target := "http://www.nicovideo.jp/watch/" + videoID
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Add("Cookie", sessionKey)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}

	u, _ := url.Parse("http://nicovideo.jp")
	tmp := jar.Cookies(u)

	return tmp[1].String(), nil
}

func GetFlv(videoID string, sessionKey string) (flv map[string]string, err error) {
	target := "http://flapi.nicovideo.jp/api/getflv?v=" + videoID
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Add("Cookie", sessionKey)

	client := &http.Client{}
	res, err := client.Do(req)

	defer res.Body.Close()

	if err != nil {
		return flv, err
	}

	byteArray, _ := ioutil.ReadAll(res.Body)
	arr := strings.Split(string(byteArray), "&")
	flv = map[string]string{}
	for i, _ := range arr {
		q := strings.Split(arr[i], "=")
		unescapedValue, _ := url.QueryUnescape(q[1])
		flv[q[0]] = unescapedValue
	}

	return flv, nil
}

func DownloadVideoSource(videoURL string, outputPath string, nicoHistory string) (err error) {
	req, _ := http.NewRequest("GET", videoURL, nil)
	req.Header.Add("Cookie", nicoHistory)

	temporaryPath := outputPath + ".nvdownload"

	// Resume download
	if stat, err := os.Stat(temporaryPath); err == nil {
		req.Header.Add("Range", "bytes="+fmt.Sprint(stat.Size())+"-")
	}

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
	dataLength, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	bar := pb.New(dataLength).SetUnits(pb.U_BYTES)
	bar.Start()

	file, err := os.OpenFile(temporaryPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	writer := io.MultiWriter(file, bar)

	io.Copy(writer, res.Body)

	// Rename when finished
	if err := os.Rename(temporaryPath, outputPath); err != nil {
		return err
	}

	return nil
}
