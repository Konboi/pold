package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var (
	Version = "0.0.1"
)

func main() {
	newApp().Run(os.Args)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "pold"
	app.Usage = "markdown based blog tool"
	app.Version = Version
	app.Author = "Konboi"
	app.Email = "ryosuke.yabuki+pold@gmail.com"
	app.Commands = Commands

	return app
}
