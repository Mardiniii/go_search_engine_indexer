package main

import (
	"fmt"

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

func crawlURL(url string) {
	// Extract links, title and description
	s := NewScraper(url)
	links := s.ScrapeLinks()
	fmt.Println("Scraped links:", len(links))
	title, description := s.MetaDataInformation()
	body := s.Body()

	// Check if the page exists
	existsLink, page := FindPage(url)

	if existsLink {
		// Update the page in database
		fmt.Println("URL:", url, "with ID:", page.ID, "already exists")
		params := map[string]interface{}{
			"title":       title,
			"description": description,
			"body":        body,
		}

		UpdatePage(page.ID, params)
	} else {
		// Create the new page in the database.
		fmt.Println("Creating new page in the databese for link:", url)
		id, _ := shortid.Generate()
		newPage := Page{
			ID:          id,
			Title:       title,
			Description: description,
			Body:        body,
			URL:         url,
		}
		CreatePage(newPage)
	}
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

	crawlURL("http://www.elcolombiano.com")
	crawlURL("http://www.makeitreal.camp")
	crawlURL("http://www.facebook.com")
	crawlURL("http://www.eltiempo.com")
	crawlURL("http://www.elespectador.com")
	crawlURL("http://www.atlnacional.com.co")

	searchForContent("WEB Developer")
	searchForContent("Capturan")
}
