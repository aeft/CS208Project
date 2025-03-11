package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type FactorizeResponse struct {
	Factors []int `json:"factors"`
}

func main() {
	dest := "localhost:8080"

	if os.Getenv("IN_DOCKER") != "" {
		dest = "nginx"
	}

	numbers := []int64{922337203685477578, 1111111}

	for _, number := range numbers {

		// Define the API URL (adjust number as needed)
		url := fmt.Sprintf("http://%s/factorize?number=%d", dest, number)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error calling API:", err)
			return
		}
		defer resp.Body.Close()

		// Read and print the response body.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var apiResponse FactorizeResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		fmt.Printf("Number: %d\nFactors: %v\n", number, apiResponse.Factors)
	}
}
