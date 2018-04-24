package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/teris-io/shortid"
)

// Page struct to store in database
type Page struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	URL         string `json:"url"`
}

var queue = make(chan string)

func crawlURL(url string) {
	// Extract links, title and description
	s := NewScraper(url)
	if s == nil {
		return
	}
	links := s.ScrapeLinks()
	title, description := s.MetaDataInformation()
	body := s.Body()

	// Check if the page exists
	existsLink, page := ExistingPage(url)

	if existsLink {
		// Update the page in database
		params := map[string]interface{}{
			"title":       title,
			"description": description,
			"body":        body,
		}

		success := UpdatePage(page.ID, params)
		if !success {
			return
		}
		fmt.Println("Page", url, "with ID", page.ID, "updated")
	} else {
		// Create the new page in the database.
		id, _ := shortid.Generate()
		newPage := Page{
			ID:          id,
			Title:       title,
			Description: description,
			Body:        body,
			URL:         url,
		}
		success := CreatePage(newPage)
		if !success {
			return
		}
		fmt.Println("Page", url, "created")
	}

	for _, link := range links {
		go func(l string) {
			queue <- l
		}(link)
	}
}

func worker(wg *sync.WaitGroup, id int) {
	for link := range queue {
		crawlURL(link)
	}
	wg.Done()
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not option provided, please specify one of the options below:\n")
		fmt.Println("1. If you want to crawl the internet:")
		fmt.Println("\tgo run *.go index CRAWLING_START_URL\n")
		fmt.Println("2. If you want to delete the pages index from elastic search:")
		fmt.Println("\tgo run *.go delete")
		return
	}

	switch args[1] {
	case "index":
		NewElasticSearchClient()
		exists := ExistsIndex(indexName)
		if !exists {
			CreateIndex(indexName)
		}

		var wg sync.WaitGroup
		noOfWorkers := 10
		start := os.Args[2]

		go func(s string) {
			queue <- s
		}(start)

		wg.Add(noOfWorkers)
		for i := 1; i <= noOfWorkers; i++ {
			go worker(&wg, i)
		}
		wg.Wait()
	case "delete":
		NewElasticSearchClient()
		DeleteIndex()
	}
}
