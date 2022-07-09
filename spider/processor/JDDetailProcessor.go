package processor

import (
	"JD-price-tracker/spider/downloader"
	"JD-price-tracker/spider/parser"
	"JD-price-tracker/spider/saver"
	"log"
)

type JDDetailProcessor struct {
}

func (p JDDetailProcessor) Run(url string) {
	// 发送第一次请求，获取基本信息
	bodyText, err := downloader.Download(url)
	if err != nil {
		log.Fatal(err)
	}
	item, err := parser.ParseJDItemBasicInfo(bodyText)
	if err != nil {
		log.Fatal(err)
	} else {
		// 若第一次请求成功，则发送第二次请求，获取价格信息和优惠券信息
		couponUrl := "https://item-soa.jd.com/getWareBusiness?skuId=" + item.SkuId
		bodyText, _ = downloader.Download(couponUrl)
		item = parser.ParseJDItemPriceInfo(bodyText, item)

		// 将 item 保存到 MongoDB
		saver.ItemSaver(item)
	}
}
