package controllers

import (
	"context"
	"fmt"
	"gastrono-go/database"
	"gastrono-go/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		var allMenus []bson.M

		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuId = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)

		if insertErr != nil {
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D

		if menu.StartDate != nil && menu.EndDate != nil {
			if !inTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				msg := "Please enter the time in correct format"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}

			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.StartDate})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.EndDate})

			if menu.Name != "" {
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}

			if menu.Category != "" {
				updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
			}

			menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})

			upsert := true

			opt := options.UpdateOptions{
				Upsert: &upsert,
			}

			result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
			if err != nil {
				msg := "Menu update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
			}

			defer cancel()
			c.JSON(http.StatusOK, result)

		}

		defer cancel()
	}
}
