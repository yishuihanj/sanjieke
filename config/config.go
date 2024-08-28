package config

import (
	"fmt"
	"sanjieke/downloader"
)

var (
	Authorization string
	Cookie        string
	ApiKey        string
	CourseId      string
	Title         string
	CourseIdMap   map[int]*Course    = make(map[int]*Course)
	CourseNameMap map[string]*Course = make(map[string]*Course)
	CourseList    []*Course          = make([]*Course, 0)
	OutDirectory  string             = "./output" // 输出目录     string             = "" // 输出目录
	VideoQuality  string             = "1080p"
)

// Course 课程数据
type Course struct {
	Name       string //名字
	NodeId     int    //课程id
	VideoList  []CourseVideo
	NextCourse *Course
}

func (c *Course) GetFileName() string {
	return fmt.Sprintf("%v%v", c.Name, downloader.TsExt)
}

// CourseVideo 课程视频源
type CourseVideo struct {
	Url             string //地址
	ResolutionRatio string // 1080p 720p
	Name            string
}
