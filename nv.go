package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var Commands = []cli.Command{
	CommandGet,
}

func main() {
	app := cli.NewApp()
	app.Name = "nv"
	app.Version = Version
	app.Usage = ""
	app.Author = "Yasuaki Uechi"
	app.Email = "uetchy@randompaper.co"
	app.Commands = Commands

	app.Run(os.Args)
}
