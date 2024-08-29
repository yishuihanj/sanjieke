package config

import (
	"fmt"
)

var Config *config

func init() {
	if Config == nil {
		Config = &config{
			CourseList:   make([]*Course, 0),
			VideoQuality: "1080p",
		}
	}
}

type config struct {
	Authorization string `yaml:"authorization"`
	Cookie        string `yaml:"cookie"`
	ApiKey        string `yaml:"sjk-apikey"`
	CourseId      string `yaml:"course_id"`
	Title         string
	CourseList    []*Course
	OutDirectory  string `yaml:"out_directory"`
	VideoQuality  string
}

// Course 课程数据
type Course struct {
	Name       string //名字
	NodeId     int    //课程id
	VideoList  []CourseVideo
	NextCourse *Course
}

func (c *Course) GetFileName() string {
	return fmt.Sprintf("%v%v", c.Name, ".ts")
}

// CourseVideo 课程视频源
type CourseVideo struct {
	Url             string //地址
	ResolutionRatio string // 1080p 720p
	Name            string
}
