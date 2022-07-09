package test

import (
	"JD-price-tracker/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("JD").Collection("item_sku")

	item := model.ItemSku{
		SkuId: "100020530545",
		Url: "https://item.jd.com/100020530544.html",
		ItemName: "Beats Flex 蓝牙无线 入耳式手机耳机 颈挂式耳机 带麦可通话 Beats 经典黑红",
		ActivityPrice: 20,
	}
	filter := bson.M{"_id": item.SkuId}
	pItem, _ := bson.Marshal(item)
	var bItem bson.M
	_ = bson.Unmarshal(pItem, &bItem)
	opts := options.Update().SetUpsert(true)

	result, _ := collection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bItem}}, opts)
	fmt.Printf("%#v", result)
}