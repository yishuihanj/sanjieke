package main

import (
	"github.com/AlecAivazis/survey/v2"
	"log"
	"sanjieke/config"
)

// 输入的配置是
func checkInputConfig() bool {
	err, authorization := input("请输入Authorization")
	if err != nil {
		log.Fatalln("input 错误", err.Error())
		return false
	}
	if authorization == "" {
		log.Fatalln("Authorization 不能为空")
		return false
	}
	config.Config.Authorization = authorization

	err, cookie := input("请输入Cookie")
	if err != nil {
		log.Fatalln("input 错误", err.Error())
		return false
	}
	if cookie == "" {
		log.Fatalln("Cookie 不能为空")
		return false
	}
	config.Config.Cookie = cookie

	err, sjkApikey := input("请输入Sjk-Apikey")
	if err != nil {
		log.Fatalln("input 错误", err.Error())
		return false
	}
	if sjkApikey == "" {
		log.Fatalln("Sjk-Apikey 不能为空")
		return false
	}
	config.Config.ApiKey = sjkApikey

	err, id := input("请输入课程Id")
	if err != nil {
		log.Fatalln("input 错误", err.Error())
		return false
	}
	if id == "" {
		log.Fatalln("课程Id 不能为空")
		return false
	}
	config.Config.CourseId = id

	err, outPut := input("请输入输出目录")
	if err != nil {
		log.Fatalln("input 错误", err.Error())
		return false
	}
	if outPut == "" {
		log.Fatalln("outPut 不能为空")
		return false
	}
	config.Config.OutDirectory = outPut
	return true
}

func input(ask string) (error, string) {
	answer := ""
	prompt := &survey.Input{
		Message: ask,
	}
	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return err, ""
	}
	return nil, answer
}
