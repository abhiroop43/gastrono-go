package routes

import {
	controller "gastrono-go/controllers"
	"github.com/gin-gonic/gin"
}

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/order_items", controller.GetOrderItems())
	incomingRoutes.GET("/order_items/:order_item_id", controller.GetOrderItem())

	incomingRoutes.POST("/order_items", controller.CreateOrderItem())
	incomingRoutes.PATCH("/order_items/:order_item_id", controller.UpdateOrderItem())
}