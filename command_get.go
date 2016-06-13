package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/spf13/viper"
	"github.com/uetchy/nv/niconico"
	"os"
	"strconv"
)

var CommandGet = cli.Command{
	Name:        "get",
	Usage:       "",
	Description: "",
	Action: func(context *cli.Context) error {
		argQuery := context.Args().Get(0)

		if argQuery == "" {
			cli.ShowCommandHelp(context, "get")
			os.Exit(1)
		}

		err := loadConfig()
		if err != nil {
			err = generateConfig()
			if err != nil {
				panic(fmt.Errorf("%s", err))
			}
		}

		email := viper.GetString("email")
		password := viper.GetString("password")
		if email == "" {
			fmt.Println("Must setup 'email' first")
			os.Exit(1)
		}
		if password == "" {
			fmt.Println("Must setup 'password' first")
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

		return nil
	},
}

func getVideo(videoID string, sessionKey string) error {
	// Create output path
	thumb, _ := niconico.GetThumbInfo(videoID)
	inv := map[string]string{
		"Title":     thumb.Title,
		"VideoID":   thumb.VideoID,
		"Extension": thumb.MovieType,
	}
	videoFilenameTmpl := "{{.Title}} [{{.VideoID}}].{{.Extension}}"
	videoDestPath := applyTemplate(videoFilenameTmpl, inv)
	commentsDestPath := videoDestPath + ".json"
	fmt.Println(videoDestPath)

	// Stop if target file already exist
	if _, err := os.Stat(videoDestPath); err == nil {
		fmt.Println("Already donwloaded: " + thumb.Title)
		return errors.New("Already donwloaded")
	}

	// Fetch meta data
	flv, _ := niconico.GetFlv(videoID, sessionKey)
	nicoHistory, _ := niconico.GetHistory(videoID, sessionKey)

	// Download video
	if err := niconico.DownloadVideoSource(flv["url"], videoDestPath, nicoHistory); err != nil {
		fmt.Println("Failed: " + thumb.Title)
		return err
	}

	// Download video comments
	commentURL := flv["ms"]
	commentLength, _ := strconv.Atoi(flv["l"])
	commentThreadID := flv["thread_id"]
	if err := niconico.DownloadVideoComments(commentURL, commentsDestPath, nicoHistory, commentThreadID, commentLength); err != nil {
		fmt.Println("Failed to fetch comments:", commentURL)
		return err
	}

	fmt.Println("Downloaded: " + thumb.Title)
	return nil
}
