package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "wiki-names/docs"
)

func SetupRouter(handler *WikiHandler) *http.Server {
	log.Println("Starting server...")
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup a results cache for 2 minutes per URI, stored in Memory (DEV and PRE) or Redis in Production
	var memoryStore *persist.MemoryStore
	var redisStore *persist.RedisStore
	if os.Getenv("APP_ENV") == "dev" {
		memoryStore = persist.NewMemoryStore(1 * time.Minute)
		router.Use(cache.CacheByRequestURI(memoryStore, 2*time.Second))
	} else {
		// In Production, speed up caching by using Redis storage instead
		redisStore = persist.NewRedisStore(redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    os.Getenv("REDISHOST"),
		}))
		router.Use(cache.CacheByRequestURI(redisStore, 2*time.Second))
	}
	router.GET("search/:name", handler.GetContent)
	router.GET("extract/:name", handler.GetExtract)
	router.GET("extract/:name/:locale", handler.GetExtract)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	handler := WikiHandler{http: &http.Server{}}
	handler.http = SetupRouter(&handler)

	log.Printf("Listening on port %v\n", handler.http.Addr)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := handler.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := handler.http.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
