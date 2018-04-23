package main

import (
	"fmt"
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

func crawlURL(wg *sync.WaitGroup, url string) {
	// Extract links, title and description
	s := NewScraper(url)
	if s == nil {
		wg.Done()
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
			wg.Done()
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
			wg.Done()
			return
		}
		fmt.Println("Page", url, "created")
	}

	for _, link := range links {
		wg.Add(1)
		go crawlURL(wg, link)
	}
	wg.Done()
}

func searchForContent(input string) {
	// Search content
	pages := SearchContent(input)

	fmt.Println(len(pages), "found it for:", input)

	// Print page
	for _, p := range pages {
		fmt.Printf("Page - Id: %s - Title: %s - Description: %s - URL: %s\n",
			p.ID,
			p.Title,
			p.Description,
			p.URL,
		)
	}

	fmt.Println()
}

func main() {
	NewElasticSearchClient()
	exists := ExistsIndex(indexName)

	if !exists {
		CreateIndex(indexName)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go crawlURL(&wg, "https://www.npmjs.com/package/elasticsearch-console")
	wg.Wait()
}
