package processor

import (
	"JD-price-tracker/spider/downloader"
	"JD-price-tracker/spider/parser"
	"JD-price-tracker/spider/saver"
)

type JDListProcessor struct {
}

func (p JDListProcessor) Run(url string) {
	// 获取列表页的信息
	bodyText, _ := downloader.Download(url)
	itemList := parser.ParseJDList(bodyText)

	// 将 itemList 保存到 MongoDB
	saver.ItemListSaver(itemList)
}

