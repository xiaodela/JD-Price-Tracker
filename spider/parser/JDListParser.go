package parser

import (
	"JD-price-tracker/model"
	"regexp"
	"strings"
	"time"
)

func ParseJDList(bodyText []byte) []model.ItemSku {
	var itemList []model.ItemSku
	// 获取 SKU ID
	re := regexp.MustCompile(`wids:'([\d,]*?)',`)
	match := re.FindSubmatch(bodyText)
	skuIds := strings.Split(string(match[1]), ",")
	for _, skuId := range skuIds {
		item := model.ItemSku{
			Site:                "www.jd.com",
			LatestDiscoveryTime: time.Now(),
			SkuId:               skuId,
			Url:				 "https://item.jd.com/" + skuId + ".html",
			Source:              "Category List",
		}
		itemList = append(itemList, item)
	}
	return itemList
}
