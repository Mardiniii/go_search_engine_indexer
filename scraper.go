package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Scraper for each website
type Scraper struct {
	url string
	doc *goquery.Document
}

// NewScraper builds a new scraper for the website
func NewScraper(u string) *Scraper {
	d, err := goquery.NewDocument(u)
	if err != nil {
		log.Fatal(err)
	}

	return &Scraper{
		url: u,
		doc: d,
	}
}

// Body returns a string with the body of the page
func (s *Scraper) Body() string {
	body := s.doc.Find("body").Text()
	// Remove leading/ending white spaces
	body = strings.TrimSpace(body)

	return body
}

// ScrapeLinks returns the links from the website
func (s *Scraper) ScrapeLinks() []string {
	links := make([]string, 0)
	var link string

	s.doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		link = ""

		linkTag := item
		href, _ := linkTag.Attr("href")

		if !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "javascript") {
			if strings.HasPrefix(href, "/") {
				link = strings.Join([]string{s.url, href}, "")
			} else {
				link = href
			}

			if link != "" {
				link = strings.TrimRight(link, "/")
				links = append(links, link)
			}
		}
	})

	return links
}

// MetaDataInformation returns the title and description from the page
func (s *Scraper) MetaDataInformation() (string, string) {
	var t string
	var d string

	t = s.doc.Find("title").Contents().Text()

	s.doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("name", "") == "description" || item.AttrOr("property", "") == "og:description" {
			d = item.AttrOr("content", "")
		}
	})

	return t, d
}
