package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Worker function which receives tasks from the jobs channel and sends results to the results channel
func worker(id int, jobs <-chan string, results chan<- string) {
	for url := range jobs {
		fmt.Printf("Worker %d started fetching %s\n", id, url)
		resp, err := http.Get(url)
		if err != nil {
			results <- fmt.Sprintf("Worker %d failed fetching %s: %v", id, url, err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			results <- fmt.Sprintf("Worker %d failed reading response body from %s: %v", id, url, err)
			continue
		}
		results <- fmt.Sprintf("Worker %d completed fetching %s, length: %d", id, url, len(body))
	}
}

func main() {
	start := time.Now()

	// List of URLs to fetch
	urls := []string{
		"https://www.golang.org",
		"https://www.google.com",
		"https://www.github.com",
	}

	// Create channels for jobs and results
	jobs := make(chan string, len(urls))
	results := make(chan string, len(urls))

	// Launch a fixed number of worker goroutines
	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// Send URLs as jobs to the jobs channel
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	// Collect results from the results channel
	for r := 0; r < len(urls); r++ {
		fmt.Println(<-results)
	}

	fmt.Printf("All tasks completed in %s\n", time.Since(start))
}