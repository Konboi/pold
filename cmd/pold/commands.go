package main

import (
	"fmt"
	//"os"

	"github.com/Konboi/pold"
	"github.com/codegangsta/cli"
)

const (
	CONFIG_PATH = "pold.yml"
)

var Commands = []cli.Command{
	commandInit,
	commandServer,
}

var commandInit = cli.Command{
	Name:        "init",
	Usage:       "Set up new blog",
	Description: "Create new blog enviroment",
	Action:      doInit,
}

var commandServer = cli.Command{
	Name:        "server",
	Usage:       "launch server",
	Description: "",
	Action:      launchServer,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "config, c", Usage: "Set config path"},
	},
}

func doInit(c *cli.Context) error {
	err := pold.Init()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func launchServer(c *cli.Context) error {
	query := c.Args().First()
	config := c.Bool("config")

	config_path := CONFIG_PATH
	if config {
		config_path = query
	}

	conf, err := pold.NewConfig(config_path)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	server := pold.NewServer(conf)
	server.Run()

	return nil
}
