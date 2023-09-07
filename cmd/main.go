package main

import (
	"context"
	"fmt"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"log"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	route "gitlab.com/naonweh-studio/bubbme-backend/api/route"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
)

// @title Tag Service API
func main() {
	fmt.Println("start application")

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Fatalf("Config file not found: %v", err)
	}

	app := bootstrap.App()
	env := app.Env
	db := app.Mysql
	cache := app.RedisCache
	timeout := time.Duration(env.ContextTimeout) * time.Second

	commonFields := map[string]interface{}{
		"userId":    "12345",
		"ipAddress": GetOutboundIP(),
	}
	ctx := context.WithValue(context.Background(), "commonFields", commonFields)

	logger.NewLoggerWrapper("zap", ctx)

	gin := gin.Default()
	gin.Use(corsMiddleware())

	route.SetupSwagger(env, timeout, db, cache, gin)
	route.Setup(env, timeout, db, cache, gin)
	route.SetupCMS(env, timeout, db, cache, gin)

	gin.Run(env.ServerAddress)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
