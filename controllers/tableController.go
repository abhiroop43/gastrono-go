package controllers

import (
	"context"
	"gastrono-go/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

var tableCollection = database.OpenCollection(database.Client, "tables")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		result, err := tableCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var allTables []bson.M

		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal()
		}
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
