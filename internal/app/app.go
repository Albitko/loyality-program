package app

import (
	"github.com/gin-gonic/gin"

	"github.com/Albitko/loyalty-program/internal/controller"
)

func Run() {
	userHandler := controller.NewUserAuthHandler()
	ordersHandler := controller.NewOrdersHandler()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/api/user/register", userHandler.Register)
	r.POST("/api/user/login", userHandler.Login)

	authorized := r.Group("/api/user")
	authorized.POST("orders", ordersHandler.CreateOrder)
	authorized.GET("orders", ordersHandler.GetOrders)
	authorized.GET("balance", ordersHandler.GetBalance)
	authorized.GET("balance/withdraw", ordersHandler.Withdraw)
	authorized.GET("withdrawals", ordersHandler.GetWithdrawn)

	r.Run(":8080")
}
