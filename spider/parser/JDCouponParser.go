package parser

import (
	"JD-price-tracker/model"
	"JD-price-tracker/util"
	"encoding/json"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func ParseJDItemPriceInfo(bodyText []byte, item model.ItemSku) model.ItemSku {
	var data map[string]interface{}
	_ = json.Unmarshal(bodyText, &data)
	if data["price"] == nil {
		return item
	}
	priceInfo := data["price"].(map[string]interface{})
	item.OriginPrice, _ = strconv.ParseFloat(priceInfo["op"].(string), 64)
	item.ActivityPrice, _ = strconv.ParseFloat(priceInfo["p"].(string), 64)

	serviceInfo := data["servicesInfoUnited"].(map[string]interface{})
	if stockInfo, ok := serviceInfo["stockInfo"].(map[string]interface{}); ok {
		if delivery, okk := stockInfo["promiseResult"].(string); okk {
			item.DeliveryInfo = delivery
		}
	}

	var storeCouponList []model.Coupon
	var platformCouponList []model.Coupon
	var priceBreakList []model.Coupon

	// 店铺优惠券
	if promotion, ok := data["promotion"].(map[string]interface{}); ok {
		if activities, ok := promotion["activity"].([]interface{}); ok {
			for _, activity := range activities {
				couponDesc := activity.(map[string]interface{})["value"].(string)
				re := regexp.MustCompile(`满(.*?)元[，可]*?减(.*?)元`)
				matches := re.FindAllSubmatch([]byte(couponDesc), -1)
				for _, match := range matches {
					// 初始化优惠券
					coupon := model.Coupon{
						Source:  "store",
						Overlap: false,
					}
					if strings.Contains(couponDesc, "每满") {
						coupon.Multi = true
					} else {
						coupon.Multi = false
					}
					coupon.Quota, _ = strconv.ParseFloat(string(match[1]), 64)
					coupon.Discount, _ = strconv.ParseFloat(string(match[2]), 64)

					// 确定属于店铺优惠券还是店铺满减
					if activity.(map[string]interface{})["text"].(string) == "满减" {
						coupon.Type = "Store Price Break"
						priceBreakList = append(priceBreakList, coupon)
					} else {
						coupon.Type = "Store Coupon"
						storeCouponList = append(storeCouponList, coupon)
					}
				}
			}
		}
	}

	// 平台优惠券
	if couponInfos, ok := data["couponInfo"].([]interface{}); ok {
		for _, couponInfo := range couponInfos {
			// 初始化优惠券
			coupon := model.Coupon{
				Source:   "Platform",
				Type:     "Platform Coupon",
				Quota:    couponInfo.(map[string]interface{})["quota"].(float64),
				Discount: couponInfo.(map[string]interface{})["discount"].(float64),
				Overlap:  couponInfo.(map[string]interface{})["overlap"].(bool),
			}
			if couponInfo.(map[string]interface{})["anotherType"].(float64) == 6 {
				coupon.Multi = true
			} else {
				coupon.Multi = false
			}
			platformCouponList = append(platformCouponList, coupon)
		}
	}
	// 将三种优惠活动均加入到 CouponInfo 字段中
	item.CouponInfo = append(append(storeCouponList, platformCouponList...), priceBreakList...)

	// 计算到手价
	var possibleBuyPrices []float64
	possibleBuyPrices = append(possibleBuyPrices, item.ActivityPrice)

	// 1. 考虑店铺券
	for _, coupon := range storeCouponList {
		if coupon.Multi {
			couponNum := math.Floor(item.ActivityPrice / coupon.Quota) // 券重复使用的次数
			possibleBuyPrices = append(possibleBuyPrices, item.ActivityPrice-couponNum*coupon.Discount)
		} else {
			if item.ActivityPrice >= coupon.Quota {
				possibleBuyPrices = append(possibleBuyPrices, item.ActivityPrice-coupon.Discount)
			}
		}
	}

	// 2. 考虑平台券
	possibleBuyPrices1 := make([]float64, len(possibleBuyPrices))
	copy(possibleBuyPrices1, possibleBuyPrices)

	for _, coupon := range platformCouponList {
		if coupon.Overlap { // 可以叠加，只要活动价大于 quota，就在考虑店铺券的基础上减去平台券的优惠
			for _, possibleBuyPrice := range possibleBuyPrices {
				if coupon.Multi {
					couponNum := math.Floor(item.ActivityPrice / coupon.Quota) // 券重复使用的次数
					possibleBuyPrices1 = append(possibleBuyPrices1, possibleBuyPrice-couponNum*coupon.Discount)
				} else {
					if item.ActivityPrice >= coupon.Quota {
						possibleBuyPrices1 = append(possibleBuyPrices1, possibleBuyPrice-coupon.Discount)
					}
				}
			}
		} else { // 不可叠加，在活动价的基础上减去平台券的优惠
			if coupon.Multi {
				couponNum := math.Floor(item.ActivityPrice / coupon.Quota) // 券重复使用的次数
				possibleBuyPrices1 = append(possibleBuyPrices1, item.ActivityPrice-couponNum*coupon.Discount)
			} else {
				if item.ActivityPrice >= coupon.Quota {
					possibleBuyPrices1 = append(possibleBuyPrices1, item.ActivityPrice-coupon.Discount)
				}
			}
		}
	}

	// 3. 考虑满减
	possibleBuyPrices2 := make([]float64, len(possibleBuyPrices1))
	copy(possibleBuyPrices2, possibleBuyPrices1)

	for _, coupon := range priceBreakList {
		for _, possibleBuyPrice := range possibleBuyPrices1 {
			if item.ActivityPrice >= coupon.Quota {
				possibleBuyPrices2 = append(possibleBuyPrices2, possibleBuyPrice-coupon.Discount)
			}
		}
	}

	item.BuyPrice = util.Min(possibleBuyPrices2)
	return item
}

