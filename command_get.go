package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/spf13/viper"
	"github.com/uetchy/nv/niconico"
	"os"
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
	filenameTmpl := "[{{.VideoID}}] {{.Title}}.{{.Extension}}"
	outputPath := applyTemplate(filenameTmpl, inv)
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
		return err
	}

	fmt.Println("Downloaded: " + thumb.Title)
	return nil
}
