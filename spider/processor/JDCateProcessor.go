package processor

import (
	"JD-price-tracker/spider/downloader"
	"JD-price-tracker/spider/parser"
	"fmt"
)

type JDCateProcessor struct {
}

func (p JDCateProcessor) Run(url string) {
	// 发送第一次请求，获取基本信息
	bodyText, _ := downloader.Download(url)
	cateList := parser.ParseJDCategory(bodyText)
	fmt.Printf("%#v\n", cateList[0])
}
