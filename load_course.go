package main

import (
	"errors"
	"fmt"
	"log"
	"path"
	"sanjieke/config"
	"sanjieke/http"
)

// 加载课程
func loadCourse() error {
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
	fmt.Println("获取课程信息Title:", config.Title)

	err = ensureDirExists(path.Join(config.OutDirectory, config.Title))
	if err != nil {
		return fmt.Errorf("创建目录失败[%v]失败[%v]", config.OutDirectory, err.Error())
	}

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
			NodeId:    s.NodeId,
			Name:      s.Name,
			VideoList: make([]config.CourseVideo, 0),
		}
		config.CourseIdMap[s.NodeId] = course
		config.CourseNameMap[s.Name] = course
		config.CourseList = append(config.CourseList, course)
		if preCourse != nil {
			course.NextCourse = preCourse
		}
		preCourse = course
	}

	////开始遍历下载视频
	for _, course := range config.CourseList {
		if FileExists(path.Join(config.OutDirectory, config.Title, course.GetFileName())) {
			fmt.Println("跳过已下载课程", course.Name)
			continue
		}
		courseResp, err := http.GetCourseNode(course.NodeId)
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
					if s.ResolutionRatio == config.VideoQuality {
						m3u8Url = s.Url
					}
				}
			}
		}
		if m3u8Url == "" {
			log.Printf("课程[%v]没有找到视频", course.Name)
			continue
		}
		//todo 下载视频
	}
	return nil
}
