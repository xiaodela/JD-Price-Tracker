package parser

import (
	"JD-price-tracker/model"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseJDItemBasicInfo(bodyText []byte) (model.ItemSku, error) {
	//fmt.Println(string(bodyText))
	// 新建一个 SKU
	item := model.ItemSku{
		Site:            "www.jd.com",
		LatestCrawlTime: time.Now(),
	}
	// 获取 skuId
	re := regexp.MustCompile(`skuid: (\d+),`)
	match := re.FindSubmatch(bodyText)
	if match == nil {
		return item, fmt.Errorf("failed to fetch item details")
	}
	skuId := string(match[1])
	item.SkuId = skuId
	item.Url = "https://item.jd.com/" + skuId + ".html"

	// 商品名称
	re = regexp.MustCompile(`name: '(.*)',`)
	match = re.FindSubmatch(bodyText)
	item.ItemName = string(match[1])

	// 店铺信息
	re = regexp.MustCompile(`shopId: ?'(\d+)',`)
	match = re.FindSubmatch(bodyText)
	shopId := string(match[1])
	item.ShopId = shopId
	item.ShopUrl = "https://mall.jd.com/index-" + shopId + ".html"

	re = regexp.MustCompile(`target="_blank" title="(.*)" clstag`)
	match = re.FindSubmatch(bodyText)
	item.ShopName = string(match[1])

	// 上架状态
	re = regexp.MustCompile(`warestatus: (\d)`)
	match = re.FindSubmatch(bodyText)
	item.Status, _ = strconv.Atoi(string(match[1]))

	// 商品类目
	re = regexp.MustCompile(`catName: \[(.*)],`)
	match = re.FindSubmatch(bodyText)
	categories := strings.Split(string(match[1]), ",")
	for _, cate := range categories {
		item.Category = append(item.Category, strings.Trim(cate, "\"")) // 去掉字符串前后的双引号
	}

	// 商品属性
	re = regexp.MustCompile(`colorSize: (.*),\s*warestatus`)
	match = re.FindSubmatch(bodyText)
	var data []map[string]interface{}
	_ = json.Unmarshal(match[1], &data)
	for _, prop := range data {
		varSkuId := strconv.FormatFloat(prop["skuId"].(float64), 'f', -1, 64)
		if varSkuId == skuId {
			item.Props = prop
		} else {
			item.RelatedSkuIds = append(item.RelatedSkuIds, varSkuId)
		}
	}
	delete(item.Props, "skuId")

	// 封面图
	re = regexp.MustCompile(`src: '(.*)',`)
	match = re.FindSubmatch(bodyText)
	item.CoverPic = "https://img1.360buyimg.com/n0/" + string(match[1])

	// 轮播图
	re = regexp.MustCompile(`imageList: (.*),`)
	match = re.FindSubmatch(bodyText)
	var imgUrls []string
	_ = json.Unmarshal(match[1], &imgUrls)
	for _, url := range imgUrls {
		item.CarouselPics = append(item.CarouselPics, "https://img1.360buyimg.com/n0/"+url)
	}

	return item, nil
}
