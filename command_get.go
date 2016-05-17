package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/spf13/viper"
	"github.com/uetchy/nv/niconico"
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

		sessionKey := niconico.GetSessionKey(email, password)
		if niconico.IsMylist(argQuery) {
			mylistID := niconico.ToMylistID(argQuery)
			mylist, _ := niconico.GetMylist(mylistID, sessionKey)
			for _, video := range mylist.List {
				getVideo(video.ID, sessionKey)
			}
		} else {
			videoID := niconico.ToVideoID(argQuery)
			getVideo(videoID, sessionKey)
		}
	},
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/nv")
	err := viper.ReadInConfig()
	return err
}

func getVideo(videoID string, sessionKey string) error {
	thumb, _ := niconico.GetThumbInfo(videoID)

	// Create output path
	rootPath := viper.GetString("root")
	filenameTmpl := "{{.ProviderURL}}/watch/{{.VideoID}}/{{.Title}}.{{.Extension}}"
	inv := map[string]string{
		"Title":       thumb.Title,
		"VideoID":     thumb.VideoID,
		"Extension":   thumb.MovieType,
		"ProviderURL": "www.nicovideo.jp",
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
	flv, _ := niconico.GetFlv(videoID, sessionKey)
	nicoHistory, _ := niconico.GetHistory(videoID, sessionKey)

	// Download video
	if err := niconico.DownloadVideoSource(flv["url"], outputPath, nicoHistory); err != nil {
		fmt.Println("Failed: " + thumb.Title)
		return errors.New("Failed download")
	}

	fmt.Println("Downloaded: " + thumb.Title)
	return nil
}
