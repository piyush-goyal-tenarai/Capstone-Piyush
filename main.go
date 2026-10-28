package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-management/config"
	"order-management/controller"
	"order-management/database"
	"order-management/metrics"
	"order-management/middleware"
	"order-management/repository"
	"order-management/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo, cfg.ProductServiceURL)
	ctrl := controller.NewOrderController(svc)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.RequestID())
	router.Use(middleware.StructuredLogger())
	router.Use(metrics.MetricsMiddleware())

	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	router.GET("/health/ready", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "db": "disconnected"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready", "db": "connected"})
	})

	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, metrics.Instance.Snapshot())
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/orders", ctrl.CreateOrder)
		v1.GET("/orders/:id", ctrl.GetOrder)
		v1.GET("/orders", ctrl.ListOrders)
		v1.PUT("/orders/:id/status", ctrl.UpdateOrderStatus)
	}

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Printf("[Server] Order Management API on http://0.0.0.0:%s\n", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n[Shutdown] Graceful shutdown initiated...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[Shutdown] Forced: %v", err)
	}
	fmt.Println("[Shutdown] Server stopped")
}
