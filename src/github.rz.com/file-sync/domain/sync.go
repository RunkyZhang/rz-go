package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.rz.com/file-sync/common"
)

func NewFileSyncer(sourceDirectoryPath string, targetDirectoryPath string) (*fileSyncer, error) {
	if !strings.HasSuffix(sourceDirectoryPath, "\\") {
		sourceDirectoryPath += "\\"
	}
	if !strings.HasSuffix(targetDirectoryPath, "\\") {
		targetDirectoryPath += "\\"
	}

	if fileInfo, err := os.Stat(sourceDirectoryPath); nil == err {
		if !fileInfo.Mode().IsDir() {
			return nil, fmt.Errorf("is not directory path(%s)", sourceDirectoryPath)
		}
	} else {
		return nil, err
	}
	if fileInfo, err := os.Stat(targetDirectoryPath); nil == err {
		if !fileInfo.Mode().IsDir() {
			return nil, fmt.Errorf("is not directory path(%s)", targetDirectoryPath)
		}
	} else {
		if !os.IsExist(err) {
			if err := common.Command.MakeDirectory(targetDirectoryPath); nil != err {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &fileSyncer{
		sourceDirectoryPath: sourceDirectoryPath,
		targetDirectoryPath: targetDirectoryPath,
		sourceFileMetas:     make(map[string]*FileMeta),
		targetFileMetas:     make(map[string]*FileMeta),
	}, nil
}

type fileSyncer struct {
	sourceDirectoryPath string
	targetDirectoryPath string

	sourceFileMetas map[string]*FileMeta
	targetFileMetas map[string]*FileMeta
}

func (f *fileSyncer) FindOut() error {
	var groupError error
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)

	go func() {
		var err error
		if f.sourceFileMetas, err = f.findOut(f.sourceDirectoryPath); nil != err {
			groupError = err
		}
		waitGroup.Done()
	}()

	go func() {
		var err error
		if f.targetFileMetas, err = f.findOut(f.targetDirectoryPath); nil != err {
			groupError = err
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()

	return groupError
}

func (f *fileSyncer) Sync(syncMode SyncMode) {
	if SyncModeCommon == syncMode {
		f.commonSync()
	} else if SyncModeClear == syncMode {

	} else {
	}
}

func (f *fileSyncer) commonSync() {
	for sourcePath, sourceFileMeta := range f.sourceFileMetas {
		targetPath := strings.Replace(sourcePath, f.sourceDirectoryPath, f.targetDirectoryPath, 1)
		targetFileMeta, exist := f.targetFileMetas[targetPath]

		if sourceFileMeta.IsDirectory {
			if exist {
				continue
			}

			if err := common.Command.MakeDirectory(targetPath); nil != err {
				fmt.Printf("failed to make directory(%s); error: %s\n", sourcePath, err.Error())
			} else {
				fmt.Printf("success to make directory(%s)\n", sourcePath)
			}
		} else {
			if exist && sourceFileMeta.ModifyTime == targetFileMeta.ModifyTime && sourceFileMeta.Size == targetFileMeta.Size {
				continue
			}

			if err := f.copyFile(sourcePath, targetPath); nil != err {
				fmt.Printf("failed to copy file from (%s) to (%s); error: %s\n", sourcePath, targetPath, err.Error())
			} else {
				fmt.Printf("success to copy file from (%s) to (%s)\n", sourcePath, targetPath)
			}
		}
	}
}

func (f *fileSyncer) findOut(path string) (map[string]*FileMeta, error) {
	fileMetas := make(map[string]*FileMeta)

	return fileMetas, filepath.Walk(path, func(path string, fileInfo os.FileInfo, err error) error {
		if nil != err {
			fmt.Printf("failed to get file(%s); err: %s\n", path, err.Error())
			return nil
		}

		fileMeta := &FileMeta{
			Path:        path,
			ModifyTime:  fileInfo.ModTime().Unix(),
			IsDirectory: fileInfo.IsDir(),
			Size:        fileInfo.Size(),
		}
		fileMetas[path] = fileMeta

		return nil
	})
}

func (f *fileSyncer) copyFile(sourceFilePath string, targetFilePath string) error {
	// make parent path
	targetFileDirectoryPath := filepath.Dir(targetFilePath)
	if _, err := os.Stat(targetFileDirectoryPath); nil != err {
		if !os.IsExist(err) {
			if err := common.Command.MakeDirectory(targetFileDirectoryPath); nil != err {
				return fmt.Errorf("failed to make parent directory(path); error: %s", err.Error())
			}
		}
	}

	return common.Command.CopyFile(sourceFilePath, targetFilePath)
}
