package handlers

import (
	"JD-price-tracker/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

type ItemHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewItemsHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *ItemHandler {
	return &ItemHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}


func (handler *ItemHandler) ListItemsHandler(c *gin.Context) {
	val, err := handler.redisClient.Get("items").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)

		items := make([]model.ItemSku, 0)
		for cur.Next(handler.ctx) {
			var item model.ItemSku
			cur.Decode(&item)
			items = append(items, item)
		}

		data, _ := json.Marshal(items)
		handler.redisClient.Set("items", string(data), 0)
		c.JSON(http.StatusOK, items)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		items := make([]model.ItemSku, 0)
		json.Unmarshal([]byte(val), &items)
		c.JSON(http.StatusOK, items)
	}
}
