package main

import (
	"flag"
	"fmt"
	"strings"

	"github.rz.com/file-sync/domain"
)

func main() {
	flag.String("sourceDirectoryPath", "D:\\Pictrues\\Photo\\Myself\\Mystery\\string\\object\\char\\bool\\DeepMystery", "source directory path")
	flag.String("targetDirectoryPath", "E:\\Pictrues\\Photo\\Myself\\Mystery\\string\\object\\char\\bool\\DeepMystery", "target directory path")
	flag.Parse()

	sourceDirectoryPath := flag.Lookup("sourceDirectoryPath").Value.String()
	targetDirectoryPath := flag.Lookup("targetDirectoryPath").Value.String()

	sourceDirectoryPath = strings.Replace(sourceDirectoryPath, "/", "\\", -1)
	targetDirectoryPath = strings.Replace(targetDirectoryPath, "/", "\\", -1)

	fileSyncer, err := domain.NewFileSyncer(sourceDirectoryPath, targetDirectoryPath)
	if nil != err {
		panic(err)
	}

	if err := fileSyncer.FindOut(); nil != err {
		panic(err)
	}

	fileSyncer.Sync(domain.SyncModeCommon)

	fmt.Println("done.")
}
