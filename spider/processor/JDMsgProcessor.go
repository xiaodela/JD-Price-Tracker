package processor

import (
	"JD-price-tracker/spider/messager"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type JDMsgProcessor struct {
}

func (p JDMsgProcessor) Run(skuId string) {
	desc := ""

	// 连接 MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	collection := client.Database("JD").Collection("item_sku")
	findOptions := options.Find()

	// 查询当前到手价和目标价格
	cur, err := collection.Find(context.TODO(), bson.M{"_id": skuId}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		err := cur.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		switch item["target_price"].(type) {
		case float64:
			targetPrice := item["target_price"].(float64)
			buyPrice := item["buy_price"].(float64)
			itemName := item["item_name"]
			url := item["url"]
			if buyPrice <= targetPrice {
				desc += fmt.Sprintf("[%s](%s) 当前到手价为 %.2f, 请及时查看。 \n",
					itemName, url, buyPrice)
			}
		}
	}

	// 发送信息到微信
	if desc != "" {
		messager.SendMessage("Good price!", desc)
	}
	// 中断连接
	cur.Close(context.TODO())
	client.Disconnect(context.TODO())
}
