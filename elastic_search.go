package main

// Elastic search client
import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/olivere/elastic"
)

const (
	indexName    = "pages"
	indexMapping = `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"page":{
				"properties":{
					"title": {
						"type":"text"
					},
					"description": {
						"type":"text"
					},
					"body": {
						"type":"text"
					},
					"url": {
						"type":"text"
					}
				}
			}
		}
	}`
)

var client *elastic.Client

// NewElasticSearchClient returns an elastic seach client
func NewElasticSearchClient() *elastic.Client {
	var err error

	// Create a new elastic client
	client, err = elastic.NewClient(
		elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Stop()

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	return client
}

// ExistsIndex checks if the given index exists or not
func ExistsIndex(i string) bool {
	// Check if index exists
	exists, err := client.IndexExists(i).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return exists
}

// CreateIndex creates a new index
func CreateIndex(i string) {
	createIndex, err := client.CreateIndex(indexName).
		Body(indexMapping).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	if !createIndex.Acknowledged {
		log.Println("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}
}

// DeleteIndex in the indexName constant
func DeleteIndex() {
	ctx := context.Background()
	deleteIndex, err := client.DeleteIndex(indexName).Do(ctx)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
	if !deleteIndex.Acknowledged {
		log.Println("DeleteIndex was not acknowledged. Check that timeout value is correct.")
	}
}

// SearchContent returns the results for a given query
func SearchContent(input string) []Page {
	var pages []Page

	ctx := context.Background()
	// Search for a page in the database using multi match query
	q := elastic.NewMultiMatchQuery(input, "title", "description", "body").
		Fuzziness("2").
		MinimumShouldMatch("2")
	result, err := client.Search().
		Index(indexName).
		Pretty(true).
		Sort("_score", false).
		Query(q).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var ttyp Page
	for _, page := range result.Each(reflect.TypeOf(ttyp)) {
		p := page.(Page)
		pages = append(pages, p)
	}

	return pages
}

// FindPage return a boolean and a page if the link is already
// stored in the database
func FindPage(link string) (bool, Page) {
	var exists bool
	var p Page

	ctx := context.Background()
	// Search for a page in the database using Term Query
	// q := elastic.NewTermQuery("url", link)
	q := elastic.NewMultiMatchQuery(link, "url")
	result, err := client.Search().
		Index(indexName).
		Query(q).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var ttyp Page
	for _, result := range result.Each(reflect.TypeOf(ttyp)) {
		page := result.(Page)
		if page.URL == link {
			exists, p = true, page
			return exists, p
		}
	}

	return exists, p
}

// CreatePage adds a new page to the database
func CreatePage(p Page) bool {
	var err error
	ctx := context.Background()

	page, err := client.Index().
		Index("pages").
		Type("page").
		Id(p.ID).
		BodyJson(p).
		Do(ctx)

	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("Indexed page %s to index %s, type %s\n", page.Id, page.Index, page.Type)
	return true
}

// UpdatePage adds a new page to the database
func UpdatePage(id string, params map[string]interface{}) bool {
	ctx := context.Background()

	update, err := client.Update().Index(indexName).Type("page").Id(id).Doc(params).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Updated page with id: ", update.Id)
	return true
}
