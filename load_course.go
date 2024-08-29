package main

import (
	"errors"
	"fmt"
	"log"
	"path"
	"sanjieke/api"
	"sanjieke/config"
	"sanjieke/downloader"
	"sanjieke/pkg/tool"
	"time"
)

// 下载课程
func downloadCourse() error {
	infoResp, err := api.GetInfo()
	if err != nil {
		return err
	}
	if infoResp.Code != 200 {
		return fmt.Errorf("code not 200,msg:%v", infoResp.Msg)
	}
	if infoResp.Data.Title == "" {
		return errors.New("title is empty")
	}
	config.Config.Title = infoResp.Data.Title
	fmt.Println("获取课程信息Title:", config.Config.Title)

	treeResp, err := api.GetTree()
	if err != nil {
		return err
	}

	if treeResp.Code != 200 {
		return errors.New("code not 200")
	}
	var preCourse *config.Course
	for _, s := range treeResp.Data.Tree {
		course := &config.Course{
			NodeId:    s.NodeId,
			Name:      s.Name,
			VideoList: make([]config.CourseVideo, 0),
		}
		config.Config.CourseList = append(config.Config.CourseList, course)
		if preCourse != nil {
			course.NextCourse = preCourse
		}
		preCourse = course
	}

	//开始遍历下载视频
	for _, course := range config.Config.CourseList {
		if tool.FileExists(path.Join(config.Config.OutDirectory, config.Config.Title, course.GetFileName())) {
			fmt.Println("跳过已下载课程", course.Name)
			continue
		}
		courseResp, err := api.GetCourseNode(course.NodeId)
		if err != nil {
			return err
		}
		if courseResp.Code != 200 {
			return fmt.Errorf("courseResp.Code err,code:%v,msg:%v", courseResp.Code, courseResp.Msg)
		}
		m3u8Url := ""
		for _, node := range courseResp.Data.Nodes {
			for _, s := range node.VideoContent.ResolutionRatioObjList {
				if m3u8Url == "" {
					m3u8Url = s.Url
				} else {
					if s.ResolutionRatio == config.Config.VideoQuality {
						m3u8Url = s.Url
					}
				}
			}
		}
		if m3u8Url == "" {
			log.Printf("课程[%v]没有找到视频", course.Name)
			continue
		}
		dl, err := downloader.NewTask(path.Join(config.Config.OutDirectory, config.Config.Title), course.GetFileName(), m3u8Url, 3)
		if err != nil {
			return err
		}
		err = dl.Start()
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}
