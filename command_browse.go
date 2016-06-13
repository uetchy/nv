package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var CommandBrowse = cli.Command{
	Name:        "browse",
	Usage:       "",
	Description: "",
	Action: func(context *cli.Context) error {
		argQuery := context.Args().Get(0)

		if argQuery == "" {
			cli.ShowCommandHelp(context, "browse")
			os.Exit(1)
		}

		return nil
	},
}

// def browse(filepath)
//   video_id = File.basename(filepath).match(/[^\w]([\w]{2}\d+)[^\w]/)[1]
//   system "open http://www.nicovideo.jp/watch/#{video_id}"
// end
