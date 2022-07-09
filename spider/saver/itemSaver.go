package saver

import (
	"JD-price-tracker/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ItemSaver(item model.ItemSku) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	//fmt.Println("Connected to MongoDB!")
	collection := client.Database("JD").Collection("item_sku")

	filter := bson.M{"_id": item.SkuId}
	pItem, _ := bson.Marshal(item)
	var bItem bson.M
	_ = bson.Unmarshal(pItem, &bItem)
	opts := options.Update().SetUpsert(true)

	result, _ := collection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bItem}}, opts)

	if result.UpsertedCount == 1 { // 插入新数据
		log.Println("Insert", result.UpsertedID)
	} else if result.ModifiedCount == 1 { // 更新数据
		log.Println("Update", item.SkuId)
	}
}
