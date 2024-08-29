// Partial reference https://github.com/oopsguy/m3u8

package downloader

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"path"
	"path/filepath"
	"sanjieke/pkg/deque"
	"sanjieke/pkg/httper"
	"sanjieke/pkg/parse"
	"sanjieke/pkg/tool"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type Downloader struct {
	lock        sync.Mutex
	queue       *deque.Deque[int]
	folder      string // output folder
	fileName    string //  xxx.ts
	finish      int32
	segLen      int
	tmpFolder   string // 临时文件目录
	result      *parse.Result
	concurrency int // 并发下载数量，不会超过切片的数量

	downloading *progressbar.ProgressBar
}

const defaultConcurrency = 2

// NewTask returns a Task instance
func NewTask(output, downloadName, url string, concurrency int) (*Downloader, error) {
	// check param
	if output == "" {
		return nil, fmt.Errorf("param output is empty")
	}
	if downloadName == "" {
		return nil, fmt.Errorf("param downloadName is empty")
	}
	if url == "" {
		return nil, fmt.Errorf("param url is empty")
	}
	result, err := parse.FromURL(url)
	if err != nil {
		return nil, err
	}

	d := &Downloader{
		folder:      output,
		fileName:    addTsIfNeeded(downloadName, ".ts"),
		result:      result,
		tmpFolder:   path.Join(output, ".m3u8temp"), //当前文件目录/temp
		concurrency: concurrency,
	}
	if d.concurrency <= 0 {
		d.concurrency = defaultConcurrency
	}

	d.segLen = len(result.M3u8.Segments) // segment count
	d.queue = genDeque(d.segLen)
	return d, nil
}

func addTsIfNeeded(original, suffix string) string {
	if !strings.HasSuffix(original, suffix) {
		return original + suffix
	}
	return original
}

func (d *Downloader) initFolder() error {
	//检测是否存在该目录，如果不存在则创建
	err, _ := tool.EnsureDirExists(d.folder)
	if err != nil {
		return fmt.Errorf("create dir:[%v] err:[%v]", d.folder, err.Error())
	}
	return nil
}

// 清理文件目录，先
func (d *Downloader) cleanTempFolder() error {
	//检测是否存在该目录，如果不存在则创建
	err, exist := tool.EnsureDirExists(d.tmpFolder)
	if err != nil {
		return fmt.Errorf("create dir:[%v] err:[%v]", d.tmpFolder, err.Error())
	}
	// 如果已经存在，则清除里面的文件
	if exist {
		err = tool.DeleteFilesInDir(d.tmpFolder)
		if err != nil {
			return err
		}
	}
	return nil
}

// Start runs downloader
func (d *Downloader) Start() error {
	fmt.Println("开始下载", d.fileName)
	// 检测是否存在folder
	if err := d.initFolder(); err != nil {
		return err
	}
	//检测文件是否已经存在
	if tool.FileExists(path.Join(d.folder, d.fileName)) {
		fmt.Printf("%v ==> had exist\n", d.fileName)
		return nil
	}
	// 检测是否存在临时目录，如果不存在则创建
	if err := d.cleanTempFolder(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	// 限制并发下载数量
	limitChan := make(chan struct{}, d.concurrency)
	d.downloading = progressbar.Default(
		int64(d.segLen),
		fmt.Sprintf("正在下载 %v", d.fileName),
	)
	for {
		//获取一个切片的index
		tsIdx, end, err := d.next()
		if err != nil {
			if end {
				break
			}
			continue
		}
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			err = d.download(tsIdx)
			if err != nil {
				// Back into the queue, retry request
				fmt.Printf("下载失败 %s\n", err.Error())
				if err := d.back(tsIdx); err != nil {
					fmt.Printf(err.Error())
				}
			}
			<-limitChan
		}(tsIdx)
		limitChan <- struct{}{}
	}
	wg.Wait()
	if err := d.merge(); err != nil {
		return err
	}
	return nil
}

// 下载切片,将切片先下载到临时目录，然后等待合并
func (d *Downloader) download(segIndex int) error {
	tsUrl := d.tsURL(segIndex)
	e, b := httper.Get(tsUrl)
	if e != nil {
		return fmt.Errorf("request %s, %s", tsUrl, e.Error())
	}
	defer b.Close()
	tmpFileName := d.segIndexTsTmpName(segIndex)
	fPath := filepath.Join(d.tmpFolder, tmpFileName)
	f, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("create file: %s, %s", tmpFileName, err.Error())
	}
	defer f.Close()

	src, err := io.ReadAll(b)
	if err != nil {
		return fmt.Errorf("read bytes: %s, %s", tsUrl, err.Error())
	}
	sf := d.result.M3u8.Segments[segIndex]
	if sf == nil {
		return fmt.Errorf("invalid segment index: %d", segIndex)
	}
	key, ok := d.result.Keys[sf.KeyIndex]
	if ok && key != "" {
		src, err = tool.AES128Decrypt(src, []byte(key),
			[]byte(d.result.M3u8.Keys[sf.KeyIndex].IV))
		if err != nil {
			return fmt.Errorf("decryt: %s, %s", tsUrl, err.Error())
		}
	}
	// https://en.wikipedia.org/wiki/MPEG_transport_stream
	// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
	// Need to remove the bytes before the SyncByte 0x47(71).
	syncByte := uint8(71) //0x47
	bLen := len(src)
	for j := 0; j < bLen; j++ {
		if src[j] == syncByte {
			src = src[j:]
			break
		}
	}
	_, err = io.Copy(io.MultiWriter(f), bytes.NewReader(src))
	if err != nil {
		return err
	}
	_ = d.downloading.Add(1)

	atomic.AddInt32(&d.finish, 1)
	return nil
}

func (d *Downloader) next() (segIndex int, end bool, err error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.queue.Len() == 0 {
		err = fmt.Errorf("queue empty")
		if d.finish == int32(d.segLen) {
			end = true
			return
		}
		// Some segment indexes are still running.
		end = false
		return
	}
	segIndex = d.queue.PopFront()
	return
}

func (d *Downloader) back(segIndex int) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if sf := d.result.M3u8.Segments[segIndex]; sf == nil {
		return fmt.Errorf("invalid segment index: %d", segIndex)
	}
	d.queue.PushBack(segIndex)
	return nil
}

// 当切片全部下载完毕后，合并切片
func (d *Downloader) merge() error {
	//检测是否有丢失的文件
	missingCount := 0
	for idx := 0; idx < d.segLen; idx++ {
		tsFilename := d.segIndexTsTmpName(idx)
		f := filepath.Join(d.tmpFolder, tsFilename)
		if _, err := os.Stat(f); err != nil {
			missingCount++
		}
	}
	if missingCount > 0 {
		fmt.Printf("[warning] %d files missing\n", missingCount)
	}
	mFilePath := filepath.Join(d.folder, d.fileName)
	mFile, err := os.Create(mFilePath)
	if err != nil {
		return fmt.Errorf("create main TS file failed：%s", err.Error())
	}
	defer mFile.Close()
	bar := progressbar.Default(int64(d.segLen), fmt.Sprintf("正在合并 %v", d.fileName))
	writer := bufio.NewWriter(mFile)
	mergedCount := 0
	for segIndex := 0; segIndex < d.segLen; segIndex++ {
		tsFilename := d.segIndexTsTmpName(segIndex)
		src, err := os.ReadFile(filepath.Join(d.tmpFolder, tsFilename))
		_, err = writer.Write(src)
		if err != nil {
			continue
		}
		mergedCount++
		_ = bar.Add(1)
	}
	_ = writer.Flush()
	//删除临时目录
	_ = os.RemoveAll(d.tmpFolder)
	if mergedCount != d.segLen {
		fmt.Printf("[warning] \n%d files merge failed", d.segLen-mergedCount)
	}
	fmt.Printf("下载完成 %v\n", d.fileName)
	return nil
}

func (d *Downloader) tsURL(segIndex int) string {
	seg := d.result.M3u8.Segments[segIndex]
	return tool.ResolveURL(d.result.URL, seg.URI)
}

// 每个切片的tmp文件名
// 下载的文件名称_segIndex.ts_tmp
func (d *Downloader) segIndexTsTmpName(segIndex int) string {
	return fmt.Sprintf("%v_%v%v", d.fileName, strconv.Itoa(segIndex), "_tmp")
}

func genDeque(len int) *deque.Deque[int] {
	d := deque.New[int](len)
	for i := 0; i < len; i++ {
		d.PushBack(i)
	}
	return d
}
