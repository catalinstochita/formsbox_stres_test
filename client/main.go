package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const serverURL = "http://127.0.0.1:6969"
const duration = 1 * time.Second

const jsonData = `
{
  "users": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "address": {
        "street": "123 Main St",
        "city": "Anytown",
        "state": "CA",
        "postalCode": "12345"
      },
      "phone": "123-456-7890",
      "isActive": true
    },
    {
      "id": 2,
      "name": "Jane Smith",
      "email": "jane.smith@example.com",
      "address": {
        "street": "456 Elm St",
        "city": "Anycity",
        "state": "NY",
        "postalCode": "67890"
      },
      "phone": "987-654-3210",
      "isActive": false
    },
    {
      "id": 3,
      "name": "Sam Johnson",
      "email": "sam.johnson@example.com",
      "address": {
        "street": "789 Maple Ave",
        "city": "Anystate",
        "state": "TX",
        "postalCode": "11223"
      },
      "phone": "567-890-1234",
      "isActive": true
    }
  ],
  "metadata": {
    "version": "1.0",
    "timestamp": "2023-09-15T10:00:00Z",
    "totalCount": 3
  }
}
`

func sendPostRequest() error {
	//data := map[string]interface{}{
	//	"name": "John",
	//	"age":  30,
	//}
	//
	//jsonData, err := json.Marshal(data)
	//if err != nil {
	//	return err
	//}

	response, err := http.Post(serverURL+"/insert", "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if string(body) != "7bf09f5b782c0fa342c3a07b4cca39d8" {
		return errors.New("MD5 FAILED")
	}

	return nil
}

func sendGetRequest() error {
	response, err := http.Get(serverURL + "/read")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if string(body) == "{\"age\":30,\"name\":\"John\"}\n" {
		return errors.New("ERROR")
	}

	return nil
}

func worker(wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()

	//ticker := time.NewTicker(1 * time.Second)
	//defer ticker.Stop()

	startTime := time.Now()

	requestCount := 0

	for {
		if time.Since(startTime) >= duration {
			break
		}
		if err := sendPostRequest(); err != nil {
			fmt.Println("Error sending POST request:", err)
		}
		if err := sendGetRequest(); err != nil {
			fmt.Println("Error sending GET request:", err)
		}
		requestCount++
	}
	results <- requestCount
}

func main() {
	var wg sync.WaitGroup
	maxThreads := 4
	results := make(chan int, maxThreads)

	for i := 0; i < maxThreads; i++ {
		wg.Add(1)
		go worker(&wg, results)
	}

	wg.Wait()
	close(results)

	totalRequests := 0
	for result := range results {
		fmt.Println(result)
		totalRequests += result
	}

	fmt.Printf("Total requests made by all threads: %d\n", totalRequests)
}
