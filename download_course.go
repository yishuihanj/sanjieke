package main

import (
	"context"
	"errors"
	"fmt"
	"path"
	"sanjieke/api"
	"sanjieke/config"
	"sanjieke/downloader"
	"sanjieke/pkg/markdown"
	"sanjieke/pkg/tool"
	"strings"
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

	resp, err := api.GetTree()
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return errors.New("get tree code not 200")
	}
	if resp.Data == nil || len(resp.Data.Trees) == 0 {
		return errors.New("tree data is empty")
	}
	for _, tree := range resp.Data.Trees {
		err = checkTreeVideo(tree)
		if err != nil {
			return err
		}
	}
	return nil
}

// 检测节点是否有视频
func checkTreeVideo(tree *api.TreeNode) error {
	if tree == nil {
		return nil
	}
	fmt.Println("检测课程", tree.Name)
	err := downloadAttribute(tree)
	if err != nil {
		return err
	}
	// 检测子节点
	for _, child := range tree.Children {
		// 先设置父节点
		child.ParentTree = tree
		err = checkTreeVideo(child)
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadAttribute(tree *api.TreeNode) error {
	if tree == nil || tree.Attribute == nil {
		fmt.Printf("课程:%v 没有可下载的视频，跳过\n", tree.Name)
		return nil
	}
	//获取课程详细信息
	resp, err := api.GetCourseNode(tree.NodeID)
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return errors.New("get course node code not 200")
	}

	if resp.Data == nil {
		return errors.New("course node data is empty")
	}

	downFolder := path.Join(config.Config.GetCourseRootFolder()) //下载的文件位置
	names := make([]string, 0)
	curTree := tree
	for curTree != nil {
		names = append(names, curTree.Name)
		curTree = curTree.ParentTree
	}
	// 倒序遍历
	for i := len(names) - 1; i >= 0; i-- {
		names[i] = tool.MakeValidFilename(names[i])
		downFolder = path.Join(downFolder, names[i])
	}
	videoIndex := 0

	var builder strings.Builder

	for _, node := range resp.Data.Nodes {
		switch node.ContentType {
		case "b-video":
			err, fileName := downVideoCourseNode(downFolder, videoIndex, tree, node)
			if err != nil {
				return err
			}
			if fileName == "" {
				continue
			}
			builder.WriteString(genVideoMarkdown(fileName))
		case "html":
			builder.WriteString(fmt.Sprintf("%v", node.HTMLContent))
		}
	}
	content := builder.String()
	if content == "" {
		return nil
	}
	_, err = markdown.Download(context.Background(), content, tree.Name, downFolder, false)
	if err != nil {
		return err
	}
	return nil
}

func genVideoMarkdown(fileName string) string {
	return fmt.Sprintf(`<video src="视频/%v" controls></video>`, fileName)
}

func downVideoCourseNode(downFolder string, videoIndex int, tree *api.TreeNode, node *api.CourseNode) (error, string) {
	// 如果内容类型不是视频，则跳过
	if node.ContentType != "b-video" {
		return nil, ""
	}
	if node.VideoContent == nil {
		return nil, ""
	}
	// 如果视频列表为空，则跳过
	if len(node.VideoContent.ResolutionRatioObjList) == 0 {
		return nil, ""
	}
	video := node.VideoContent.ResolutionRatioObjList[0]
	// 获取符合的视频之类和视频连接
	for _, v := range node.VideoContent.ResolutionRatioObjList {
		if v.ResolutionRatio == config.Config.VideoQuality {
			video = v
		}
	}
	//如果下载的地址为空，则跳过
	if video.URL == "" {
		return nil, ""
	}
	exs := ""
	if videoIndex > 0 {
		exs = fmt.Sprintf("_%v", videoIndex)
	}
	fileName := tool.MakeValidFilename(fmt.Sprintf("%v%v", tree.Name, exs))
	downFolder = path.Join(downFolder, "视频")
	task, err := downloader.NewTask(downFolder, fileName, video.URL, 2)
	if err != nil {
		return err, ""
	}
	err = task.Start()
	if err != nil {
		return err, ""
	}
	return nil, task.GetFileName()
}
