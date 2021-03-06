package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
	"github.com/uetchy/nv/niconico"
	"github.com/urfave/cli"
)

var CommandGet = cli.Command{
	Name:        "get",
	Usage:       "",
	Description: "",
	Flags: []cli.Flag{
		cli.BoolTFlag{
			Name:  "with-comments, c",
			Usage: "fetch comments",
		},
	},
	Action: func(context *cli.Context) error {
		argQuery := context.Args().Get(0)
		withComments := context.Bool("with-comments")
		if argQuery == "" {
			return cli.NewExitError("No argument specified", 1)
		}

		err := loadConfig()
		if err != nil {
			err = generateConfig()
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}

		email := viper.GetString("email")
		password := viper.GetString("password")
		if email == "" {
			return cli.NewExitError("Must setup 'email' first", 1)
		}
		if password == "" {
			return cli.NewExitError("Must setup 'password' first", 1)
		}

		err, sessionKey := niconico.GetSessionKey(email, password)
		if err != nil {
			return cli.NewExitError("Failed to retrieve session key", 1)
		}

		if niconico.IsMylist(argQuery) {
			mylistID := niconico.ToMylistID(argQuery)
			mylist, _ := niconico.GetMylist(mylistID, sessionKey)
			for _, video := range mylist.List {
				getVideo(video.ID, sessionKey, withComments)
			}
		} else {
			videoID := niconico.ToVideoID(argQuery)
			err := getVideo(videoID, sessionKey, withComments)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}

		return nil
	},
}

func getVideo(videoID string, sessionKey string, withComments bool) error {
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
		return errors.New("Already donwloaded" + thumb.Title)
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
	if withComments {
		commentURL := flv["ms"]
		commentLength, _ := strconv.Atoi(flv["l"])
		commentThreadID := flv["thread_id"]
		if err := niconico.DownloadVideoComments(commentURL, commentsDestPath, nicoHistory, commentThreadID, commentLength); err != nil {
			return err
		}
	}

	return nil
}
