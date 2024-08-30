package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sanjieke/config"
	"sanjieke/pkg"
)

func main() {
	app := cli.App{
		Name:    "下载器",
		Flags:   Flags(),
		Version: "1.0.0",
		Action:  mainLoop,
	}
	if err := app.Run(os.Args); err != nil {
		log.Println("!!程序运行错误 ===> ", err.Error())
		return
	}

}
func mainLoop(c *cli.Context) error {
	if cmdFlags.Ffmpeg != "" {
		config.Config.FfmpegPath = cmdFlags.Ffmpeg
	}

	if cmdFlags.Config != "" {
		err := pkg.YamlReader(cmdFlags.Config, config.Config)
		if err != nil {
			return fmt.Errorf("读取配置文件错误:%v", err.Error())
		}
	} else {
		if !checkInputConfig() {
			return fmt.Errorf("请检查输入参数")
		}
	}

	if config.Config.FfmpegPath == "" {
		return fmt.Errorf("请检查ffmpeg路径")
	}

	err := downloadCourse()
	if err != nil {
		return fmt.Errorf("下载课程错误:%v", err.Error())
	}

	return nil
}
