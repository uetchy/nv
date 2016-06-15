package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
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

		if niconico.IsMylist(argQuery) {
			mylistID := niconico.ToMylistID(argQuery)
			mylist, _ := niconico.GetMylist(mylistID, sessionKey)
			for _, video := range mylist.List {
				getVideo(video.ID, sessionKey, withComments)
			}
		} else {
			videoID := niconico.ToVideoID(argQuery)
			getVideo(videoID, sessionKey, withComments)
		}
		fmt.Println(videoID)

		return nil
	},
}

// def info(ptr)
//   config = Nv::Config.new(Nv::CONFIG_PATH)
//   config.verify_for_authentication!('info')
//
//   nico = Niconico::Base.new.sign_in(config.email, config.password)
//
//   if mylist?(ptr)
//     mylist = nico.mylist(ptr)
//
//     puts "Title : #{mylist.title}"
//     puts "Desc  : #{mylist.description}"
//
//     mylist.items.each_with_index do |item, i|
//       puts "   #{i + 1}. #{item.title}"
//     end
//   else
//     video = nico.video(ptr)
//
//     puts video.title
//     puts '=' * 40
//     puts video.description
//     puts '=' * 40
//     puts "URL: #{video.watch_url}"
//   end
// end
