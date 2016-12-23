package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var SessionId = cli.StringFlag{
	Name:  "app, a",
	Usage: "ID of session",
}

func ShowHelpError(c *cli.Context, msg string, cmd string) {
	if cmd != "" {
		cmd = " " + cmd
	}
	if c != nil {
		fmt.Fprintf(c.App.Writer, "%v%v: %v See '%v%v -h (or --help)'\n",
			c.App.Name, cmd, msg, c.App.Name, cmd)
	} else {
		fmt.Printf("%v%v: %v See '%v%v -h (or --help)'\n",
			c.App.Name, cmd, msg, c.App.Name, cmd)
	}
	os.Exit(1)
}

func GetSessionID(c *cli.Context, mandatory bool, cmd string) string {
	var name string = ""
	name = c.String("a")
	if name == "" && mandatory {
		ShowHelpError(c, "msg", cmd)
		os.Exit(1)
	}
	return name
}
