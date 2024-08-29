package main

import "github.com/urfave/cli/v2"

type flags struct {
	Config string
}

var cmdFlags flags

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "c",
			Usage:       "当使用-c参数时，使用配置文件",
			Value:       "",
			Destination: &cmdFlags.Config,
		},
	}
}
