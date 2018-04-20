package main

import (
	"fmt"
	"log"
	// Go Query Plugin for scraping
	"github.com/PuerkitoBio/goquery"
)

func linkScrape() {
	doc, err := goquery.NewDocument("http://localhost:3000/")
	if err != nil {
		log.Fatal(err)
	}

	// use CSS selector found with the browser inspector
	// for each, use index and item
	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		linkText := linkTag.Text()
		fmt.Printf("Link #%d: '%s' - '%s'\n", index, linkText, link)
	})

	var metaDescription string
	var pageTitle string

	// use CSS selector found with the browser inspector
	// for each, use index and item
	pageTitle = doc.Find("title").Contents().Text()

	doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("name", "") == "description" || item.AttrOr("property", "") == "og:description" {
			metaDescription = item.AttrOr("content", "")
		}
	})
	fmt.Printf("Page Title: '%s'\n", pageTitle)
	fmt.Printf("Meta Description: '%s'\n", metaDescription)
}

func main() {
	linkScrape()
}
