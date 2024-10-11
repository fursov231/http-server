package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	for {
		errorCount := 0
		response, err := http.Get("http://srv.msk01.gigacorp.local")
		if err != nil {
			errorCount += 1
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			errorCount += 1
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			errorCount += 1
			checkErrorLimitExceeded(errorCount)
			return
		}

		content := string(body)
		values := strings.Split(content, ",")

		var ramLimit = 0
		var spaceLimit = 0
		var throughputLimit = 0
		for i, strVal := range values {
			value, err := strconv.Atoi(strVal)

			if err != nil {
				errorCount += 1
				checkErrorLimitExceeded(errorCount)
			}

			if i == 0 {
				if value > 30 {
					fmt.Printf("Load Average is too high: %d\n", value)
				}
			}

			if i == 1 {
				ramLimit = value
			}

			if i == 2 {
				memoryUsed := (float64(value) / float64(ramLimit)) * 100
				if memoryUsed > 80 {
					fmt.Printf("Memory usage too high: %d%%\n", int(memoryUsed))
				}
			}

			if i == 3 {
				spaceLimit = value
			}

			if i == 4 {
				spaceUsagePercentage := float64(value) / float64(spaceLimit) * 100
				if spaceUsagePercentage > 90 {
					leftSpace := (spaceLimit - value) / 1_048_576
					fmt.Printf("Free disk space is too low: %d Mb left\n", leftSpace)
				}
			}

			if i == 5 {
				throughputLimit = value
			}

			if i == 6 {
				throughputPercentage := float64(value) / float64(throughputLimit) * 100
				availableThroughput := (throughputLimit - value) / 1_000_000
				if throughputPercentage > 90 {
					fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", availableThroughput)
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func checkErrorLimitExceeded(tryCount int) {
	if tryCount >= 3 {
		fmt.Println("Unable to fetch server statistic")
	}
}
