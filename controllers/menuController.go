package controllers

import (
	"context"
	"gastrono-go/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, nil)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
