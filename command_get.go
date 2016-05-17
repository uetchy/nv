package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"text/template"
)

var CommandGet = cli.Command{
	Name:        "get",
	Usage:       "",
	Description: "",
	Action: func(context *cli.Context) {
		argQuery := context.Args().Get(0)

		if argQuery == "" {
			cli.ShowCommandHelp(context, "get")
			os.Exit(1)
		}

		err := loadConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		email := viper.GetString("email")
		password := viper.GetString("password")
		root := viper.GetString("root")
		if email == "" {
			fmt.Println("Must setup 'email' first")
			os.Exit(1)
		}
		if password == "" {
			fmt.Println("Must setup 'password' first")
			os.Exit(1)
		}
		if root == "" {
			fmt.Println("Must setup 'root' first")
			os.Exit(1)
		}

		sessionKey := getSessionKey(email, password)
		if isMylist(argQuery) {
			mylistId := toMylistId(argQuery)
			mylist, _ := getMylist(mylistId, sessionKey)
			for _, video := range mylist.List {
				getVideo(video.Id, sessionKey)
			}
		} else {
			videoId := toVideoId(argQuery)
			getVideo(videoId, sessionKey)
		}
	},
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/nv")
	err := viper.ReadInConfig()
	return err
}

func getVideo(videoId string, sessionKey string) error {
	thumb, _ := getThumbInfo(videoId)

	// Create output path
	rootPath := viper.GetString("root")
	filenameTmpl := "{{.ProviderUrl}}/watch/{{.VideoId}}/{{.Title}}.{{.Extension}}"
	inv := map[string]string{
		"Title":       thumb.Title,
		"VideoId":     thumb.VideoId,
		"Extension":   thumb.MovieType,
		"ProviderUrl": "www.nicovideo.jp",
	}
	t := template.New("outputPath")
	template.Must(t.Parse(filenameTmpl))
	var buf bytes.Buffer
	t.Execute(&buf, inv)
	outputPath := filepath.Join(rootPath, buf.String())
	fmt.Println(outputPath)

	// Stop if target file already exist
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Println("Already donwloaded: " + thumb.Title)
		return errors.New("Already donwloaded")
	}

	// Fetch meta data
	flv, _ := getFlv(videoId, sessionKey)
	nicoHistory, _ := getNicoHistory(videoId, sessionKey)

	// Download video
	if err := downloadVideoSource(flv["url"], outputPath, nicoHistory); err != nil {
		fmt.Println("Failed: " + thumb.Title)
		return errors.New("Failed download")
	}

	fmt.Println("Downloaded: " + thumb.Title)
	return nil
}
