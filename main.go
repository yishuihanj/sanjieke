package main

import (
	"errors"
	"fmt"
	"log"
	"sanjieke/chromedper"
	"sanjieke/config"
	"sanjieke/http"
	"time"
)

func main() {
	err := chromedper.InitChrome()
	if err != nil {
		log.Println("初始化谷歌浏览器失败", err.Error())
		return
	}
	err = chromedper.Login()
	if err != nil {
		log.Println("登录失败", err.Error())
		return
	}
	isVip, err := chromedper.CheckVip()
	if err != nil {
		log.Println("检查是否为VIP失败", err.Error())
		return
	}
	if !isVip {
		log.Println("非VIP用户，不能下载视频")
		return
	}
	log.Println("是VIP用户，开始进行下一步")

	err = loadCourse(8005775)
	if err != nil {
		log.Println("loadCourse 错误", err.Error())
		return
	}
}

// 加载课程
func loadCourse(id int32) error {
	time.Sleep(time.Second * 1)
	//先清除缓存
	chromedper.ClearCache()
	config.ClearCache()
	err := chromedper.NavigateCourse(id)
	if err != nil {
		return err
	}

	if !config.CheckConfigValid() {
		return errors.New("配置文件无效")
	}

	infoResp, err := http.GetInfo()
	if err != nil {
		return err
	}
	if infoResp.Code != 200 {
		return errors.New("code not 200")
	}
	if infoResp.Data.Title == "" {
		return errors.New("title is empty")
	}
	config.Title = infoResp.Data.Title

	treeResp, err := http.GetTree()
	if err != nil {
		return err
	}

	if treeResp.Code != 200 {
		return errors.New("code not 200")
	}
	var preCourse *config.Course
	for _, s := range treeResp.Data.Tree {
		course := &config.Course{
			Id:        s.NodeId,
			Name:      s.Name,
			VideoList: make([]config.CourseVideo, 0),
		}
		config.CourseIdMap[s.NodeId] = course
		config.CourseNameMap[s.Name] = course
		if preCourse != nil {
			course.NextCourse = preCourse
		}
		preCourse = course
	}
	fmt.Println(config.CourseNameMap, config.CourseIdMap)
	return nil
}
