package main

import (
  "net/http"
  "io/ioutil"
  "regexp"
  "encoding/json"
)

type MylistVideoThumbnailStyle struct {
  OffsetX int
  OffsetY int
  Width   int
}

type MylistVideo struct {
  ThumbnailUrl      string
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
  Id                string
  CreateTime        int
  ThreadUpdateTime  string
  MylistComment     string
}

type Mylist struct {
  Name                 string
  Description          string
  UserId               int
  UserNickname         string
  DefaultSort          int
  List                 []MylistVideo
  IsWatchingThisMylist bool
  IsWatchingCountFull  bool
  Status               string
}

func isMylist(query string) bool {
  if m, _ := regexp.MatchString("mylist", query); !m {
    return false
  }
  return true
}

func toMylistId(query string) string {
  re, _ := regexp.Compile("\\d+")
  one := re.Find([]byte(query))

  return string(one)
}

func getMylist(mylistId string, sessionKey string) (Mylist, error) {
  target := "http://riapi.nicovideo.jp/api/watch/mylistvideo?id="
  target += mylistId
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
