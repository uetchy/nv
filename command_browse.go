package main

import (
	"os/exec"

	"github.com/urfave/cli"
)

var CommandBrowse = cli.Command{
	Name:        "browse",
	Usage:       "",
	Description: "",
	Action: func(context *cli.Context) error {
		argQuery := context.Args().Get(0)
		if argQuery == "" {
			return cli.NewExitError("No argument specified", 1)
		}

		videoID := fetchVideoIDFromFilename(argQuery)
		videoURL := "http://www.nicovideo.jp/watch/" + videoID
		if err := exec.Command("open", videoURL).Run(); err != nil {
			return err
		}

		return nil
	},
}
