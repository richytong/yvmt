# Youtube Video Metadata Tools

## Getting Started

<b><i>Advisory</i></b> - This project runs on go1.11.1 and takes advantage of [Go Modules](https://github.com/golang/go/wiki/Modules). Please clone the repo into a directory outside of your `$GOPATH` or set `GO111MODULE=on`.

`git clone git@github.com:richytong/yvmt` - clones the repo

`go build fetch.go ii.go util.go` - builds `fetch`

`go build search.go ii.go util.go` - builds `search`

## Fetch
Fetch and store youtube videos based on a search query string

Usage: `./fetch`

Flags
  - query: `--query "your youtube query here"`
  - maxResults: `--max-results 60`
  - requestsPerMin: `--requests-per-min 60`

## Search
Search from all the fetched youtube videos using an inverted index

Usage: `./search`

Flags
  - query: `--query "you local query here"`

## Exploration Guide
`fetch.go` - source code for `fetch`

`search.go` - source code for `search`

`ii.go` - defines the InvertedIndex type among other types as well as relevant methods

`util.go` - contains the extractLowerAlphanumericFields definition
