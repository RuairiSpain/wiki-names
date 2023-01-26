package app

import (
	"log"
	"net/http"
	"os"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"wiki-names/controllers"
	"wiki-names/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(address string) *http.Server {
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
	router.GET("search/:name", wiki_controller.GetContentSummary)
	router.GET("extract/:name", wiki_controller.GetExtract)
	router.GET("extract/:name/:locale", wiki_controller.GetExtract)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &http.Server{
		Addr:    address,
		Handler: router,
	}
}
