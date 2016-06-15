package niconico

import (
	"bytes"
	"encoding/json"
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
	Embeddable    bool   `xml:"thumb>embeddable"`
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

// <chat thread="1387572949" no="5" vpos="15095" date="1387575544" mail="184" user_id="kwpNfx7idiioqlDzTDbG49P-drU" anonymity="1" leaf="2">綺麗だな~</chat>
type Comment struct {
	ThreadID  int    `xml:"thread,attr" json:"thread"`
	No        int    `xml:"no,attr" json:"no"`
	Vpos      int    `xml:"vpos,attr" json:"vpos"`
	Date      string `xml:"date,attr" json:"date"`
	Mail      string `xml:"mail,attr" json:"mail"`
	UserID    string `xml:"user_id,attr" json:"user_id"`
	Anonymity bool   `xml:"anonymity,attr" json:"anonymity"`
	Leaf      int    `xml:"leaf,attr" json:"leaf"`
	Body      string `xml:",chardata" json:"body"`
}

// <thread resultcode="0" thread="1387572949" last_res="256" ticket="0x50cac000" revision="1" server_time="1465812057"/>
type Thread struct {
	ResultCode int    `xml:"resultcode,attr" json:"result_code"`
	ID         int    `xml:"thread,attr" json:"id"`
	LastRes    int    `xml:"last_res,attr" json:"last_res"`
	Ticket     string `xml:"ticket,attr" json:"ticket"`
	Revision   int    `xml:"revision,attr" json:"revision"`
	ServerTime int    `xml:"server_time,attr" json:"server_time"`
}

// <leaf thread="1387572949" count="49"/>
// <leaf thread="1387572949" leaf="1" count="14"/>
type Leaf struct {
	ThreadID int `xml:"thread,attr" json:"thread_id"`
	Number   int `xml:"leaf,attr" json:"number"`
	Count    int `xml:"count,attr" json:"count"`
}

// <view_counter video="32767" id="sm22495319" mylist="55"/>
type ViewCounter struct {
	Video  string `xml:"video,attr" json:"video"`
	ID     string `xml:"id,attr" json:"id"`
	Mylist int    `xml:"mylist,attr" json:"mylist"`
}

// <global_num_res thread="1387572949" num_res="256"/>
type GlobalNumberRes struct {
	NumRes int `xml:"num_res,attr" json:"num_res"`
}

type Packet struct {
	Threads         []Thread        `xml:"thread" json:"threads"`
	Leaves          []Leaf          `xml:"leaf" json:"leaves"`
	ViewCounter     ViewCounter     `xml:"view_counter" json:"view_counter"`
	GlobalNumberRes GlobalNumberRes `xml:"global_num_res" json:"global_num_res"`
	Comments        []Comment       `xml:"chat" json:"comments"`
}

func ToVideoID(query string) string {
	re, _ := regexp.Compile("[a-z]{2}?\\d+")
	one := re.FindString(query)

	return one
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

// GET http://flapi.nicovideo.jp/api/getwaybackkey?thread=1345476375
// => waybackkey=1417346808.E9d0LUF9gvFvt3Rrf5TP91Pa0LA
//
// POST http://msg.nicovideo.jp/53/api/
// <packet><thread thread="1345476375" version="20090904" user_id="1501297" scores="1" nicoru="1" with_global="1"/><thread_leaves thread="1345476375" user_id="1501297" scores="1" nicoru="1">0-14:100,1000</thread_leaves></packet>

// commentURL flv.ms
// threadID flv.thread_id
// length flv.l
func DownloadVideoComments(commentURL string, outputPath string, nicoHistory string, threadID string, length int) (err error) {
	length = int(length / 60)
	reqBody := `<packet><thread thread="` + threadID + `" version="20090904" scores="1" nicoru="1" with_global="1"/><thread_leaves thread="` + threadID + `" scores="1" nicoru="1">0-` + fmt.Sprint(length) + `:10</thread_leaves></packet>`
	req, _ := http.NewRequest("POST", commentURL, bytes.NewBuffer([]byte(reqBody)))
	req.Header.Add("Cookie", nicoHistory)

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()

	byteArray, _ := ioutil.ReadAll(res.Body)
	packet := new(Packet)
	err = xml.Unmarshal(byteArray, packet)
	if err != nil {
		return err
	}

	data, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	temporaryPath := outputPath + ".nvdownload"
	file, err := os.OpenFile(temporaryPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	file.Write(data)

	// Rename when finished
	if err := os.Rename(temporaryPath, outputPath); err != nil {
		return err
	}

	return nil
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
