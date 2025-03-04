package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	// router.Use(ResponseTimeMiddleware())

	router.GET("/factorize", factorizeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}

// func ResponseTimeMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		c.Next()
// 		duration := time.Since(start)
// 		c.Writer.Header().Set("X-Response-Time", duration.String())
// 	}
// }

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
