package main

import "github.com/urfave/cli/v2"

type flags struct {
	Config string
	Ffmpeg string
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

		&cli.StringFlag{
			Name:        "f",
			Usage:       "当使用-f参数时，配置ffmpeg的路径",
			Value:       "./ffmpeg.exe",
			Destination: &cmdFlags.Ffmpeg,
		},
	}
}
