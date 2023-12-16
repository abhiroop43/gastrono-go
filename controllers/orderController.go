package controllers

import (
	"context"
	"fmt"
	"gastrono-go/database"
	"gastrono-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var orderCollection = database.OpenCollection(database.Client, "order")
var tableCollection = database.OpenCollection(database.Client, "table")

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table
		var order models.Order

		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()

			return
		}

		if order.TableId != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
			defer cancel()

			if err != nil {
				msg := fmt.Sprintf("Table with id %s not found", *order.TableId)
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}
		}

		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.OrderId = order.ID.Hex()

		result, insertErr := orderCollection.InsertOne(ctx, order)
		defer cancel()

		if insertErr != nil {
			msg := fmt.Sprintf("Order was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table
		var order models.Order

		var updateObj primitive.D

		orderId := c.Param("order_id")
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		if order.TableId != nil {
			err := menuCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("Table with id %s not found", *order.TableId)
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.TableId})
		}

		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})

		upsert := true

		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(ctx, filter, bson.D{
			{"%st", updateObj},
		}, &opt)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("Order was not updated")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func OrderItemOrderCreator(order models.Order) string {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()

	_, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer cancel()

	return order.OrderId
}
