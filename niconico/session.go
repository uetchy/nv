package niconico

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"errors"
)

func GetSessionKey(mail string, password string) (error, string) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	if mail == "" || password == "" {
		return errors.New("No email or password given"), ""
	}

	res, err := client.PostForm("https://secure.nicovideo.jp/secure/login?site=niconico",
		url.Values{
			"mail":     {mail},
			"password": {password},
		})
	defer res.Body.Close()

	if err != nil {
		return err, ""
	}

	u, _ := url.Parse("http://nicovideo.jp")
	sessionKey := jar.Cookies(u)[1].String()

	return nil, sessionKey
}
