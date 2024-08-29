package downloader

import "testing"

var outPut = "./AAA"

func TestDownloader(t *testing.T) {
	downloader, err := NewTask(outPut, "我的祖国.ts", "xxxx", 3)
	if err != nil {
		t.Error(err)
		return
	}
	err = downloader.Start()
	if err != nil {
		t.Error(err)
		return
	}
}
