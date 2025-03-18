package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type FactorizeResponse struct {
	Factors []int `json:"factors"`
}

func discoverService(consulAddr, service string) (string, error) {
	url := fmt.Sprintf("http://%s/v1/catalog/service/%s", consulAddr, service)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to query consul: %w", err)
	}
	defer resp.Body.Close()

	var services []struct {
		ServiceAddress string `json:"ServiceAddress"`
		ServicePort    int    `json:"ServicePort"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return "", fmt.Errorf("failed to decode consul response: %w", err)
	}
	if len(services) == 0 {
		return "", fmt.Errorf("no instances for service %q found", service)
	}
	// Randomly pick one instance
	chosen := services[rand.Intn(len(services))]
	return fmt.Sprintf("%s:%d", chosen.ServiceAddress, chosen.ServicePort), nil
}

func StartFactorizationRequests(ctx context.Context) error {
	// Use Consul to dynamically discover the address for nginx.
	consulAddr := os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		consulAddr = "consul:8500"
	}

	nginxDest, err := discoverService(consulAddr, "nginx")
	if err != nil {
		// fmt.Println("Error discovering nginx service:", err)
		// nginxDest = "localhost:8080"
		panic(fmt.Sprintf("Error discovering nginx service: %s", err))
	}

	fmt.Println("Discovered nginx instance:", nginxDest)

	file, err := os.Open("test_numbers.txt")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var numbersFromFile []int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		number, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			fmt.Printf("Invalid number in file: %s\n", line)
			continue
		}
		numbersFromFile = append(numbersFromFile, number)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	seed := int64(42) // Default seed
	if seedStr := os.Getenv("TEST_SEED"); seedStr != "" {
		if parsedSeed, err := strconv.ParseInt(seedStr, 10, 64); err == nil {
			seed = parsedSeed
		} else {
			fmt.Println("Invalid TEST_SEED value, using default seed 42")
		}
	}

	testGoroutines := int32(1)
	if testGoroutinesStr := os.Getenv("TEST_GOROUTINES"); testGoroutinesStr != "" {
		if parsedTestGoroutines, err := strconv.ParseInt(testGoroutinesStr, 10, 64); err == nil {
			testGoroutines = int32(parsedTestGoroutines)
		} else {
			fmt.Println("Invalid TEST_GOROUTINES value, using default value 1")
		}
	}

	g, ctx := errgroup.WithContext(ctx)

	for t := 0; t < int(testGoroutines); t++ {
		g.Go(func() error {
			numbers := make([]int64, len(numbersFromFile))
			_ = copy(numbers, numbersFromFile)

			rand.New(rand.NewSource(seed+int64(t))).Shuffle(len(numbers), func(i, j int) {
				numbers[i], numbers[j] = numbers[j], numbers[i]
			})

			for i, number := range numbers {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				url := fmt.Sprintf("http://%s/factorize?number=%d", nginxDest, number)
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("Error calling API:", err)
					continue
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
					continue
				}

				var apiResponse FactorizeResponse
				if err := json.Unmarshal(body, &apiResponse); err != nil {
					fmt.Println("Error parsing JSON:", err)
					continue
				}

				fmt.Printf("(%d) Number: %d\nFactors: %v\n", i+1, number, apiResponse.Factors)
			}
			return nil
		})
	}

	fmt.Printf("Finish factorization\n")

	return g.Wait()
}

func main() {
	router := gin.New()
	router.GET("/start", func(c *gin.Context) {
		if err := StartFactorizationRequests(c.Request.Context()); err != nil {
			fmt.Printf("Error: %s\n", err)
			c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf(("Error: %s"), err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Finish factorization"})
	})
	router.Run(":8090")
}
