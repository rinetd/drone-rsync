package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/rinetd/drone-rsync/cmd"
	"github.com/urfave/cli"
)

var Version string = fmt.Sprintf("1.0.2")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "Drone-rsync"
	app.Usage = "Rsync to Remote Hosts"
	app.Copyright = "Copyright (c) 2017 rinetd"
	app.Authors = []cli.Author{
		{
			Name:  "rinetd",
			Email: "sdlylshl@gmail.com",
		},
	}
	app.Action = cmd.Run
	app.Version = Version
	app.Flags = []cli.Flag{

		cli.StringSliceFlag{
			Name:   "hosts,host,H",
			Usage:  "connect to host",
			EnvVar: "PLUGIN_HOSTS",
		},
		cli.StringFlag{
			Name:   "port,p",
			Usage:  "connect to port ",
			EnvVar: "PLUGIN_PORT",
			Value:  "22",
		},
		cli.StringFlag{
			Name:   "username,user,u",
			Usage:  "connect as user ",
			EnvVar: "PLUGIN_USERNAME,PLUGIN_USER",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "ssh-key,key",
			Usage:  "private ssh key",
			EnvVar: "PLUGIN_SSH_KEY,PLUGIN_KEY",
		},
		cli.StringFlag{
			Name:   "password,P",
			Usage:  "user password",
			EnvVar: "PLUGIN_PASSWORD",
		},
		cli.StringFlag{
			Name:   "source",
			Usage:  "source",
			EnvVar: "PLUGIN_SOURCE",
			Value:  ".",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "target",
			EnvVar: "PLUGIN_TARGET",
		},
		cli.StringSliceFlag{
			Name:   "args",
			Usage:  "args commands",
			EnvVar: "PLUGIN_ARGS",
		},
		cli.BoolFlag{
			Name:   "recursive",
			Usage:  "recursive mode",
			EnvVar: "PLUGIN_RECURSIVE",
		},
		cli.BoolFlag{
			Name:   "delete",
			Usage:  "delete mode",
			EnvVar: "PLUGIN_DELETE",
		},
		cli.BoolFlag{
			Name:   "sync",
			Usage:  "sync mode",
			EnvVar: "PLUGIN_SYNC",
		},
		cli.StringFlag{
			Name:   "chmod",
			Usage:  "chmod commands",
			EnvVar: "PLUGIN_CHMOD",
		},
		cli.StringFlag{
			Name:   "verbose",
			Usage:  "verbose -q -v -vv -vvv",
			EnvVar: "PLUGIN_VERBOSE",
		},
		cli.StringFlag{
			Name:   "chown",
			Usage:  "chown commands",
			EnvVar: "PLUGIN_CHOWN",
		},
		cli.StringSliceFlag{
			Name:   "include",
			Usage:  "include commands",
			EnvVar: "PLUGIN_INCLUDE",
		},
		cli.StringSliceFlag{
			Name:   "exclude",
			Usage:  "exclude commands",
			EnvVar: "PLUGIN_EXCLUDE",
		},
		cli.StringSliceFlag{
			Name:   "filter",
			Usage:  "filter commands",
			EnvVar: "PLUGIN_FILTER",
		},
		cli.StringSliceFlag{
			Name:   "export",
			Usage:  "export envirment to remote hosts",
			EnvVar: "PLUGIN_EXPORT",
		},
		cli.StringFlag{
			Name:   "scripts,script,s",
			Usage:  "execute commands",
			EnvVar: "PLUGIN_SCRIPTS,PLUGIN_SCRIPT",
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("drone-rsync error: ", err)
		os.Exit(1)
	}
}
