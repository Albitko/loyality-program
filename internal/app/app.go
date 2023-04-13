package app

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/Albitko/loyalty-program/internal/controller"
	"github.com/Albitko/loyalty-program/internal/repo"
	"github.com/Albitko/loyalty-program/internal/usecase"
	"github.com/Albitko/loyalty-program/internal/workers"
)

func Run() {
	// Implement config with ENV and FLAG

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := repo.NewRepository(ctx, "postgresql://localhost:5432/postgres")
	defer storage.Close()

	workers.InitWorkers(ctx, storage, "https://test-service.com")

	userAuthenticator := usecase.NewAuthenticator(storage)
	ordersProcessor := usecase.NewOrdersProcessor(storage)
	balanceProcessor := usecase.NewBalanceProcessor(storage)

	userHandler := controller.NewUserAuthHandler(userAuthenticator)
	ordersHandler := controller.NewOrdersHandler(ordersProcessor)
	balanceHandler := controller.NewBalanceHandler(balanceProcessor)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/api/user/register", userHandler.Register)
	r.POST("/api/user/login", userHandler.Login)

	authorized := r.Group("/api/user")
	// authorized.Use() JWT middleware
	authorized.POST("orders", ordersHandler.CreateOrder)
	authorized.GET("orders", ordersHandler.GetOrders)
	authorized.GET("balance", balanceHandler.GetBalance)
	authorized.GET("balance/withdraw", balanceHandler.Withdraw)
	authorized.GET("withdrawals", balanceHandler.GetWithdrawn)

	r.Run(":8080")
}
