package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const FactorizePath = "/factorize"

var (
	activeConnections int64
	// connectionGauge is a Prometheus gauge for active connections.
	connectionGauge = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Current number of active connections.",
		},
		func() float64 {
			return float64(atomic.LoadInt64(&activeConnections))
		},
	)

	// executionTimeHistogram records the execution time of /factorize API requests.
	executionTimeHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "execution_time_seconds",
		Help:    fmt.Sprintf("Execution time of %s API requests.", FactorizePath),
		Buckets: prometheus.DefBuckets,
	})
)

func init() {
	prometheus.MustRegister(connectionGauge)
	prometheus.MustRegister(executionTimeHistogram)
}

func ConnectionCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == FactorizePath {
			atomic.AddInt64(&activeConnections, 1)
			defer atomic.AddInt64(&activeConnections, -1)
		}
		c.Next()
	}
}

func ExecutionTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == FactorizePath {
			startTime := time.Now()
			c.Next()
			duration := time.Since(startTime).Seconds()
			executionTimeHistogram.Observe(duration)
		} else {
			c.Next()
		}
	}
}

func main() {
	router := gin.New()

	router.Use(ConnectionCountMiddleware())
	router.Use(ExecutionTimeMiddleware())

	router.GET(FactorizePath, factorizeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Expose Prometheus metrics endpoint.
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run(":" + port)
}

func factorizeHandler(c *gin.Context) {
	numStr := c.Query("number")
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil || num < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide a valid number (>= 2)"})
		return
	}
	factors := primeFactors(num)
	c.JSON(http.StatusOK, gin.H{
		"factors": factors,
	})
}

func primeFactors(n int64) []int64 {
	var factors []int64
	for i := int64(2); i*i <= n; i++ {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}
