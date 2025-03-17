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

func main() {
	dest := "localhost:8080"

	if os.Getenv("IN_DOCKER") != "" {
		dest = "nginx"
	}

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
				url := fmt.Sprintf("http://%s/factorize?number=%d", dest, number)
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

	fmt.Printf("Finish factorization.\n")

	// keep the client running
	for {
		time.Sleep(10 * time.Second)
	}
}
