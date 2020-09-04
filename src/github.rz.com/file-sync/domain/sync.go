package domain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func NewFileSyncer(sourceFolderPath string, targetFolderPath string) (*fileSyncer, error) {
	if fileInfo, err := os.Stat(sourceFolderPath); nil == err {
		if !fileInfo.Mode().IsDir() {
			return nil, fmt.Errorf("is not folder path(%s)", sourceFolderPath)
		}
	}
	if fileInfo, err := os.Stat(targetFolderPath); nil == err {
		if !fileInfo.Mode().IsDir() {
			return nil, fmt.Errorf("is not folder path(%s)", targetFolderPath)
		}
	}
	if !strings.HasSuffix(sourceFolderPath, "\\") {
		sourceFolderPath += "\\"
	}
	if !strings.HasSuffix(targetFolderPath, "\\") {
		targetFolderPath += "\\"
	}

	return &fileSyncer{
		sourceFolderPath: sourceFolderPath,
		targetFolderPath: targetFolderPath,
	}, nil
}

type fileSyncer struct {
	sourceFolderPath string
	targetFolderPath string

	sourceFileMetas map[string]*FileMeta
	targetFileMetas map[string]*FileMeta
}

func (f *fileSyncer) FindOut() error {
	var groupError error
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)

	go func() {
		var err error
		if f.sourceFileMetas, err = f.findOut(f.sourceFolderPath); nil != err {
			groupError = err
		}
		waitGroup.Done()
	}()

	go func() {
		var err error
		if f.targetFileMetas, err = f.findOut(f.targetFolderPath); nil != err {
			groupError = err
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()

	return groupError
}

func (f *fileSyncer) findOut(path string) (map[string]*FileMeta, error) {
	var fileMetas map[string]*FileMeta

	return fileMetas, filepath.Walk(path, func(path string, fileInfo os.FileInfo, err error) error {
		if nil != err {
			fmt.Printf("failed to get file(%s); err: %s\n\t", path, err.Error())
			return nil
		}

		fileMeta := &FileMeta{
			Path:       path,
			ModifyTime: fileInfo.ModTime().Unix(),
		}
		fileMetas[path] = fileMeta

		return nil
	})
}

func (f *fileSyncer) copyFile(sourceFilePath string) error {
	sourceFilePath = strings.Replace(sourceFilePath, f.sourceFolderPath, f.targetFolderPath, 1)
	return exec.Command("cmd", "/C", "copy", sourceFilePath, sourceFilePath, "/y").Run()
}

//func (f *fileSyncer) newFolder(sourceFolderPath string) error {
//	sourceFilePath = strings.Replace(sourceFilePath, f.sourceFolderPath, f.targetFolderPath, 1)
//	return exec.Command("cmd", "/C", "copy", sourceFilePath, sourceFilePath, "/y").Run()
//}