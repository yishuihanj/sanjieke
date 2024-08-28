package chromedper

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
	"sanjieke/config"
	"strings"
	"time"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func InitChrome() error {
	// 创建一个带有显示界面的 chromedp 上下文
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("mute-audio", false), // 关闭声音
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-sandbox", false),
		chromedp.Flag("disable-setuid-sandbox", false),
		chromedp.Flag("disable-infobars", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("disable-dev-shm-usage", false),
		chromedp.Flag("disable-background-networking", false),
	)
	ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
	// 创建浏览器上下文
	ctx, cancel = chromedp.NewContext(ctx)
	// 设置超时
	ctx, cancel = context.WithCancel(ctx)
	// 启用网络事件
	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		return err
	}
	listenerTarget()
	return nil
}

// 模拟登录
func Login() error {
	// 执行 chromedp 任务
	url := "https://passport.sanjieke.cn/account/sign_in?channel=youshangjiao&oauth_callback=&v=20230216"
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`input[type="text"][placeholder="输入手机号"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[type="text"][placeholder="输入手机号"]`, "15732158959", chromedp.ByQuery),

		chromedp.WaitVisible(`input[type="password"][placeholder="输入密码"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[type="password"][placeholder="输入密码"]`, "hanjie123", chromedp.ByQuery),

		chromedp.WaitVisible(`input[name="denglu_page_denglu"]`, chromedp.ByQuery), // 等待元素可见
		chromedp.Click(`input[name="denglu_page_denglu"]`, chromedp.ByQuery),       // 点击按钮
	)
	if err != nil {
		return err
	}
	return nil
}

func CheckVip() (bool, error) {
	var vipIconExists bool
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(`span.header-container_user-container_module_avator`, chromedp.ByQuery),
		// 检查 <i class="vip-icon"> 是否存在于 <span> 元素中
		chromedp.Evaluate(`document.querySelector('span.header-container_user-container_module_avator i.vip-icon') !== null`, &vipIconExists),
	)
	if err != nil {
		return false, err
	}

	return vipIconExists, nil
}

// 导航到课程，获取课程id，auth和cookie等
func NavigateCourse(courseId int32) error {
	url := fmt.Sprintf("https://www.sanjieke.cn/course/detail/sjk/%v", courseId)
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div.button-container`, chromedp.ByQuery),
	)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 3)

	// 点击播放按钮
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector('.video-icon').dispatchEvent(new MouseEvent("click",{bubbles:true}));`, nil),
	)
	if err != nil {
		return err
	}

	log.Println("延时3秒检测是否存在[去学习]按钮")
	//延时3秒
	time.Sleep(time.Second * 3)

	//判断去学习是否存在
	var letStudyButton bool
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector('button.el-button.el-button--primary span') !== null`, &letStudyButton),
	)
	if err != nil {
		return err
	}
	//如果存在去学习，需要点击
	if letStudyButton {
		// 点击播放按钮
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.querySelector('button.el-button.el-button--primary').dispatchEvent(new MouseEvent("click",{bubbles:true}));`, nil),
		)
		if err != nil {
			return err
		}
	}
	<-ch
	return nil
}

var authorization string
var cookie string
var apiKey string
var studyId string
var ch = make(chan struct{})

func listenerTarget() {
	// 监听网络请求事件
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if event, ok := ev.(*network.EventRequestWillBeSentExtraInfo); ok {
			if authorization != "" && cookie != "" && apiKey != "" {
				return
			}
			h := event.Headers

			_path, ok := h[":path"]
			if !ok {
				return
			}
			pathStr := _path.(string)
			if !strings.Contains(pathStr,
				"content/tree") {
				return
			}

			if studyId == "" {
				re := regexp.MustCompile(`/b-side/api/web/study/0/(\d+)/content/tree`)
				matches := re.FindStringSubmatch(pathStr)
				if len(matches) > 1 {
					studyId = matches[1]
				}
			}
			checkHeader(event.Headers)
		}
	})
}

func checkHeader(headers network.Headers) {
	a, ok := headers["authorization"]
	if ok && authorization == "" {
		authorization = a.(string)
	}
	c, ok := headers["cookie"]
	if ok && cookie == "" {
		cookie = c.(string)
	}
	api, ok := headers["sjk-apikey"]
	if ok && apiKey == "" {
		apiKey = api.(string)
	}
	if authorization != "" && cookie != "" && apiKey != "" {
		config.Authorization = authorization
		config.Cookie = cookie
		config.ApiKey = apiKey
		config.StudyId = studyId
		ch <- struct{}{}
	}
}

func ClearCache() {
	authorization = ""
	cookie = ""
	studyId = ""
	apiKey = ""
	ch = make(chan struct{})
}
