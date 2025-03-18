package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
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

func main() {
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
		fmt.Println("Error opening file:", err)
		return
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
		fmt.Println("Error reading file:", err)
		return
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

	wg := sync.WaitGroup{}

	for t := 0; t < int(testGoroutines); t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			numbers := make([]int64, len(numbersFromFile))
			_ = copy(numbers, numbersFromFile)

			rand.New(rand.NewSource(seed+int64(t))).Shuffle(len(numbers), func(i, j int) {
				numbers[i], numbers[j] = numbers[j], numbers[i]
			})

			for i, number := range numbers {
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
		}()
	}

	wg.Wait()

	fmt.Printf("Finish factorization.\n")

	// keep the client running
	for {
		time.Sleep(10 * time.Second)
	}
}
