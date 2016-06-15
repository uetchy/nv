package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/uetchy/nv/niconico"
	"github.com/urfave/cli"
)

var CommandInfo = cli.Command{
	Name:        "info",
	Usage:       "",
	Description: "",
	Action: func(context *cli.Context) error {
		argQuery := context.Args().Get(0)
		if argQuery == "" {
			return cli.NewExitError("No argument specified", 1)
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
				fmt.Println(video.Title)
				fmt.Println(video.DescriptionShort, "\n")
			}
		} else {
			videoID := niconico.ToVideoID(argQuery)
			video, _ := niconico.GetThumbInfo(videoID)
			fmt.Println(video.Title)
			fmt.Println(video.Description)
			fmt.Println(video.ViewCounter, "watches", video.CommentNum, "comments", video.MylistCounter, "listed")
			fmt.Println(video.WatchURL)
		}

		return nil
	},
}
