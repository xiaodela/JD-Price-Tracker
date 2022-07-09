package model

import "time"

type ItemSku struct {
	Site                string                 `bson:"site,omitempty" json:"site"`
	Url                 string                 `bson:"url,omitempty" json:"url"`
	SkuId               string                 `bson:"_id,omitempty" json:"sku_id"`
	ItemName            string                 `bson:"item_name,omitempty" json:"item_name"`
	ShopId              string                 `bson:"shop_id,omitempty" json:"shop_id"`
	ShopUrl             string                 `bson:"shop_url,omitempty" json:"shop_url"`
	ShopName            string                 `bson:"shop_name,omitempty" json:"shop_name"`
	RelatedSkuIds       []string               `bson:"related_sku_ids,omitempty" json:"related_sku_ids"`
	Status              int                    `bson:"status,omitempty" json:"status"` // 1 表示上架，2 表示下架
	Category            []string               `bson:"category,omitempty" json:"category"`
	Props               map[string]interface{} `bson:"props,omitempty" json:"props"`
	CoverPic            string                 `bson:"cover_pic,omitempty" json:"cover_pic"`
	CarouselPics        []string               `bson:"carousel_pics,omitempty" json:"carousel_pics"`
	DeliveryInfo        string                 `bson:"delivery_info,omitempty" json:"delivery_info"`
	CouponInfo          []Coupon               `bson:"coupon_info,omitempty" json:"coupon_info"`
	OriginPrice         float64                `bson:"original_price,omitempty" json:"original_price"`
	ActivityPrice       float64                `bson:"activity_price,omitempty" json:"activity_price"`
	TargetPrice         float64                `bson:"target_price,omitempty" json:"target_price"`
	BuyPrice            float64                `bson:"buy_price,omitempty" json:"buy_price"`
	LatestCrawlTime     time.Time              `bson:"latest_crawl_time,omitempty" json:"latest_crawl_time"`
	LatestDiscoveryTime time.Time              `bson:"latest_discovery_time,omitempty" json:"latest_discovery_time"`
	Source              string                 `bson:"source,omitempty" json:"source"`
}

type Coupon struct {
	Type     string
	Source   string
	Quota    float64
	Discount float64
	Overlap  bool
	Multi    bool
}
