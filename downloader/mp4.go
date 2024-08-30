package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sanjieke/config"
)

// mergeTSFiles 使用 ffmpeg 合并 .ts 文件
func (d *Downloader) mergeTSFiles(listFileName, outputFileName string) error {
	cmd := exec.Command(config.Config.FfmpegPath, "-loglevel", "quiet", "-f", "concat", "-safe", "0", "-i", listFileName, "-c", "copy", outputFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	_ = d.downloading.Add(1)
	_ = d.downloading.Finish()
	return nil
}

func (d *Downloader) genMp4() error {
	files := make([]string, 0)
	for idx := 0; idx < d.segLen; idx++ {
		tsFilename := d.segIndexTsTmpName(idx)
		f := filepath.Join(d.tmpFolder, tsFilename)
		files = append(files, f)
	}
	fileListName := filepath.Join(d.tmpFolder, "filelist.txt")
	err := createFileList(fileListName, files)
	if err != nil {
		return err
	}
	defer os.Remove(fileListName) // 合并后删除临时文件
	// 2. 使用 ffmpeg 合并文件
	mFilePath := filepath.Join(d.folder, d.fileName)
	err = d.mergeTSFiles(fileListName, mFilePath)
	if err != nil {
		return err
	}
	return nil
}

func createFileList(fileName string, files []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, f := range files {
		_, err := file.WriteString(fmt.Sprintf("file '%s'\n", f))
		if err != nil {
			return err
		}
	}

	return nil
}
