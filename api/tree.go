package api

import (
	"encoding/json"
	"fmt"
	"sanjieke/config"
	"sanjieke/pkg/httper"
)

func GetTree() (*TreeResp, error) {
	// 创建一个新的 HTTP 请求
	url := fmt.Sprintf("https://web-api.sanjieke.cn/b-side/api/web/study/0/%v/content/tree", config.Config.CourseId)
	err, all := httper.HttpGet(url)
	if err != nil {
		return nil, err
	}
	treeResp := new(TreeResp)

	err = json.Unmarshal(all, treeResp)
	if err != nil {
		return nil, err
	}
	return treeResp, nil
}

type TreeResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *Data  `json:"data"`
}

// 视频属性
type Attribute struct {
	Type         string   `json:"type"` //只有 video 才下载
	ContentTypes []string `json:"contentTypes"`
}
type TreeNode struct {
	NodeID     int         `json:"nodeId"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Attribute  *Attribute  `json:"attribute"` //课程属性
	Children   []*TreeNode `json:"children"`  //一节课程有多个小课程
	ParentTree *TreeNode   // 父节点
}
type Data struct {
	Trees []*TreeNode `json:"tree"`
}
