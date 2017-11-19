package cmd

import (
	"github.com/rinetd/drone-rsync/rsync"
	"github.com/urfave/cli"
)

func Run(c *cli.Context) error {
	p := &rsync.Plugin{
		Config: rsync.Config{
			Hosts:     c.StringSlice("hosts"),
			Port:      c.Int("port"),
			User:      c.String("user"),
			Key:       c.String("ssh-key"),
			Password:  c.String("password"),
			Source:    c.String("source"),
			Target:    c.String("target"),
			Recursive: c.Bool("recursive"),
			Delete:    c.Bool("delete"),
			Sync:      c.Bool("sync"),
			Chmod:     c.String("chmod"),
			Chown:     c.String("chown"),
			Include:   c.StringSlice("include"),
			Exclude:   c.StringSlice("exclude"),
			Filter:    c.StringSlice("filter"),
			Script:    c.String("script"),
		},
	}
	// fmt.Println("%v", p)

	p.Exec()
	return nil
}
