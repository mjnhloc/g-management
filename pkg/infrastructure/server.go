package infrastructure

import (
	"g-management/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method and path",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func NewServer(db *gorm.DB, redisClient *services.RedisClient) *gin.Engine {
	// Create logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create router with default middleware
	router := gin.New()

	// Add Redis client to context
	router.Use(func(c *gin.Context) {
		c.Set("redisClient", redisClient)
		c.Next()
	})

	// Configure custom middleware
	ConfigureMiddleware(router, logger)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// Welcome endpoint
	router.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Welcome to the g-management API",
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
		})
	})

	return router
}
