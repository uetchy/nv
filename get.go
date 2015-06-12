package main

import (
  "os"
  "fmt"
  "net/http"
  "net/http/cookiejar"
  "net/url"
  "io/ioutil"
  "strings"

  "github.com/codegangsta/cli"
  "github.com/cheggaaa/pb"
)

func getSessionKey(mail string, password string) string {
  jar, _ := cookiejar.New(nil)
  client := &http.Client{
    Jar: jar,
  }

  res, err := client.PostForm("https://secure.nicovideo.jp/secure/login?site=niconico",
    url.Values{
      "mail": {mail},
      "password": {password},
    })

  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  defer res.Body.Close()

  u, _ := url.Parse("http://nicovideo.jp")
  sessionKey := jar.Cookies(u)[1].String()

  return sessionKey
}

func getFlv(videoId string, sessionKey string) map[string]string {
  u := "http://flapi.nicovideo.jp/api/getflv?v="
  u += videoId

  req, _ := http.NewRequest("GET", u, nil)
  req.Header.Add("Cookie", sessionKey)

  client := &http.Client{}
  res, _ := client.Do(req)

  defer res.Body.Close()

  byteArray, _ := ioutil.ReadAll(res.Body)
  arr := strings.Split(string(byteArray), "&")
  flv := map[string]string{}
  for i, _ := range arr {
    q := strings.Split(arr[i], "=")
    unescapedValue, _ := url.QueryUnescape(q[1])
    flv[q[0]] = unescapedValue
  }

  return flv
}

func getNicoHistory(videoId string, sessionKey string) string {
  target := "http://www.nicovideo.jp/watch/"
  target += videoId

  req, _ := http.NewRequest("GET", target, nil)
  req.Header.Add("Cookie", sessionKey)

  jar, _ := cookiejar.New(nil)
  client := &http.Client{Jar: jar}
  res, _ := client.Do(req)

  defer res.Body.Close()

  u, _ := url.Parse("http://nicovideo.jp")
  tmp := jar.Cookies(u)
  return tmp[1].String()
}

func getVideo(videoUrl string, nicoHistory string) {
  req, _ := http.NewRequest("GET", videoUrl, nil)
  req.Header.Add("Cookie", nicoHistory)
  client := &http.Client{}
  res, _ := client.Do(req)

  defer res.Body.Close()

  fmt.Println(res.Header)
}

var Get = cli.Command{
  Name: "get",
  Usage: "",
  Description: "",
  Action: func(context *cli.Context) {
    argQuery := context.Args().Get(0)

    if argQuery == "" {
      cli.ShowCommandHelp(context, "get")
      os.Exit(1)
    }

    videoId := "sm18032359"

    sessionKey := getSessionKey("whsque@gmail.com", "in17nicoxs@+")
    flv := getFlv(videoId, sessionKey)
    fmt.Println(len(flv))
    nicoHistory := getNicoHistory(videoId, sessionKey)
    getVideo(flv["url"], nicoHistory)
  },
}
