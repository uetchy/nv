package main

import (
	"os"
	"github.com/codegangsta/cli"
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
