# Go Search Engine Indexer

This repository contains a search engine indexer used to crawling the internet from a given website as a start point. This code is built to be consumed for the search engine client in this [repository](https://github.com/Mardiniii/go_search_engine_client).

## Usage

Clone the project to your local a machine and run any of the commands below into the project directory:

1. Crawl the internet

```go
go build
./go_search_engine_indexer index START_CRAWLING_URL
```

Or, you can run the next command:

```go
go run *.go index START_CRAWLING_URL
```

2. Delete the current elastic search index

```go
go build
./go_search_engine_indexer delete
```

Or, you can run the next command:

```go
go run *.go delete
```

## Contributions
Feel free to make any comment, pull request, code review, shared post, fork or feedback. Everything is welcome.

## License

This project is licensed under the **MIT License**.

## Authors

**Sebastian Zapata Mardini** - [GitHub profile](https://github.com/Mardiniii)
