package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aungmyozaw92/go-graphql/cmd"
	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/directives"
	"github.com/aungmyozaw92/go-graphql/graph"
	"github.com/aungmyozaw92/go-graphql/middlewares"
	"github.com/aungmyozaw92/go-graphql/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ravilushqa/otelgqlgen"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

const defaultPort = "8080"

var tracer = otel.Tracer("go-graphql")

type Cache struct {
	client redis.UniversalClient
	ttl    time.Duration
}

const apqPrefix = "apq:"

func getRedisClient(redisAddress string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	return client
}

func NewCache(redisAddress string, ttl time.Duration) (*Cache, error) {

	client := getRedisClient(redisAddress)

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("could not create cache: %w", err)
	}

	return &Cache{client: client, ttl: ttl}, nil
}

// func (c *Cache) Add(ctx context.Context, key string, value interface{}) {
// 	c.client.Set(context.Background(), apqPrefix+key, value, c.ttl)
// }

// func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool) {
// 	s, err := c.client.Get(context.Background(), apqPrefix+key).Result()
// 	if err != nil {
// 		return struct{}{}, false
// 	}
// 	return s, true
// }

func (c *Cache) Add(ctx context.Context, key string, value string) {
	c.client.Set(context.Background(), apqPrefix+key, value, c.ttl)
}

func (c *Cache) Get(ctx context.Context, key string) (string, bool) {
	s, err := c.client.Get(context.Background(), apqPrefix+key).Result()
	if err != nil {
		return "", false
	}
	return s, true
}

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file

	cache, err := NewCache(os.Getenv("REDIS_ADDRESS"), 24*time.Hour)
	if err != nil {
		panic("cannot create APQ redis cache")
	}

	c := graph.Config{Resolvers: &graph.Resolver{
		Tracer: tracer,
	}}
	c.Directives.Auth = directives.Auth

	h := handler.NewDefaultServer(graph.NewExecutableSchema(c))
	h.Use(otelgqlgen.Middleware())
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{
		MaxMemory:     32 << 20, // 32 MB
		MaxUploadSize: 50 << 20, // 50 MB
	})
	h.Use(extension.AutomaticPersistedQuery{Cache: cache})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}


func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

    cmd.Execute()

	logger := config.GetLogger()

	// Connect to Database
	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	models.MigrateTable()
	// Initialize Gin router.
	r := gin.New()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	corsConfig.AddAllowHeaders("token", "Origin", "Content-Type", "Authorization")
	corsConfig.AddExposeHeaders("Content-Length")
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))
	r.Use(middlewares.AuthMiddleware())
	r.Use(middlewares.LoaderMiddleware())
	r.Use(customErrorLogger(logger))
	r.Use(gin.Recovery())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	
	r.NoRoute(customNotFoundHandler)
	r.Run(":" + port)

	logger.WithFields(logrus.Fields{
		"info": "Connection Established",
	}).Info("connect to http://localhost:", port, "/ for GraphQL playground")
	log.Println("Server started successfully")
}

// customErrorLogger is a custom Gin middleware that logs only errors
func customErrorLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only log when there are errors
		if len(c.Errors) > 0 {
			logger.Error(c.Errors.String())
		}
	}
}

func customNotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
}
