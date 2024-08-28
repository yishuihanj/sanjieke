package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sanjieke/config"
	"time"
)

func GetTree() (*TreeResp, error) {
	// 创建一个新的 HTTP 请求
	url := fmt.Sprintf("https://web-api.sanjieke.cn/b-side/api/web/study/0/%v/content/tree", config.StudyId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	setNormalHeader(req.Header)
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
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
	Data struct {
		Tree []struct {
			NodeId    int    `json:"nodeId"`
			Name      string `json:"name"`
			Type      string `json:"type"`
			Attribute struct {
				IsFinish      int      `json:"isFinish"`
				Obligatory    int      `json:"obligatory"`
				Type          string   `json:"type"`
				IsCurrentNode int      `json:"isCurrentNode"`
				ContentTypes  []string `json:"contentTypes"`
			} `json:"attribute"`
			Children      interface{} `json:"children"`
			UnLockTime    interface{} `json:"unLockTime"`
			IsUnLock      int         `json:"isUnLock"`
			VideoDuration int         `json:"videoDuration"`
			ExamCount     int         `json:"examCount"`
		} `json:"tree"`
	} `json:"data"`
}
