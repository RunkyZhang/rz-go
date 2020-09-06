package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	words := map[string]string{
		"pal":       "老兄",
		"list":      "目录",
		"town":      "城镇",
		"computer":  "电脑",
		"crazy":     "疯狂的",
		"climate":   "气候",
		"ocean":     "海洋",
		"company":   "公司",
		"kitchen":   "",
		"printed":   "尖的",
		"waiter":    "服务生",
		"kill":     "杀死",
		"matter":    "物质",
		"person":    "人",
		"scenery":   "风景",
		"romantic":  "浪漫的",
		"mind":      "介意",
		"artist":    "艺术家",
		"expensive": "昂贵的",
		"reason":    "理由",
		"calendar":  "日历",
		"hate":      "讨厌",
		"skill":     "技术",
		"corner":    "角落",
		"household": "家庭",
		"popular":   "受欢迎的",
		"amazing":   "神奇的",
		"terrible":  "可怕的",
		"medical":   "医学",
		"clam":      "平静",
		"patient":   "病人",
		"expert":    "专家",
		"frankness": "坦白，真诚",
		"protect":   "保护",
		"lead":      "领导",
		"repair":    "修理",
		"hardly":    "几乎不",
		"treat":     "治疗，款待",
		"business":  "商务，业务",
		"goods":     "货物",
		"neighbor":  "邻居",
		"kid":       "小孩",
		"trip":      "旅行",
		"match":     "比赛",
		"suppose":   "假设",
		"careless":  "粗心",
	}

	content := generateContent(words, 100)

	file, err := os.Create("D:/english.doc")
	if err != nil {
		panic(fmt.Sprintf("io.Create(); error: %s", err.Error()))
	}
	_, err = io.Copy(file, bytes.NewReader([]byte(content)))
	if err != nil {
		panic(fmt.Sprintf("io.Copy(); error: %s", err.Error()))
	}
	fmt.Println("done")
}

func generateContent(words map[string]string, englishPer int) string {
	englishMaxLength := 0
	chineseMaxLength := 0
	for english, chinese := range words {
		if len(english) > englishMaxLength {
			englishMaxLength = len(english)
		}

		if len(chinese) > chineseMaxLength {
			chineseMaxLength = len(chinese)
		}
	}
	chineseMaxLength = chineseMaxLength / 3 * 2

	content := ""
	index := 1
	for english, chinese := range words {
		value := random.Intn(100)

		blank := ""
		underline := ""
		if value < englishPer {
			blankLength := englishMaxLength - len(english)
			for i := 0; i < blankLength; i++ {
				blank += " "
			}
			for i := 0; i < chineseMaxLength; i++ {
				underline += "_"
			}

			content += fmt.Sprintf("%d. %s", index, english+blank+":"+underline)
		} else {
			for i := 0; i < englishMaxLength; i++ {
				underline += "_"
			}

			blankLength := chineseMaxLength - len(chinese)
			for i := 0; i < blankLength; i++ {
				blank += " "
			}

			content += fmt.Sprintf("%d. %s", index, underline+":"+chinese+blank)
		}

		index++
	}

	return content
}
