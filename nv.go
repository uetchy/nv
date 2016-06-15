package main

import (
	"github.com/urfave/cli"
	"os"
)

var Commands = []cli.Command{
	CommandGet,
	CommandInfo,
	CommandBrowse,
}

var Version string = "HEAD"

func main() {
	app := cli.NewApp()
	app.Name = "nv"
	app.Version = Version
	app.Usage = "nv get [URL | NAME]"
	app.Author = "Yasuaki Uechi"
	app.Email = "uetchy@randompaper.co"
	app.Commands = Commands

	app.Run(os.Args)
}
