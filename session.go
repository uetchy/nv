package main

import (
  "os"
  "fmt"
  "net/http"
  "net/http/cookiejar"
  "net/url"
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
