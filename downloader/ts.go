package downloader

import (
	"bufio"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
)

func (d *Downloader) genTs() error {
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
	if mergedCount != d.segLen {
		fmt.Printf("[warning] \n%d files merge failed", d.segLen-mergedCount)
	}
	return nil
}
