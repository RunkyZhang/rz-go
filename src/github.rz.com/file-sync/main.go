package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println(filepath.Dir("D:\\opt\\report.exe"))

	flag.String("sourcePath", "D:/codes", "source path")
	flag.String("targetPath", "D:/test", "target path")
	flag.Parse()

	sourcePath := flag.Lookup("sourcePath").Value.String()
	targetPath := flag.Lookup("targetPath").Value.String()

	sourcePath = strings.Replace(sourcePath, "/", "\\", -1)
	targetPath = strings.Replace(targetPath, "/", "\\", -1)

	fmt.Println("done.")
}
