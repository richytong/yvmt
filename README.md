# Youtube Video Metadata Tools

## Getting Started
`go build fetch.go ii.go util.go` - builds `fetch`

`go build search.go ii.go util.go` - builds `search`

## Fetch
Fetch and store youtube videos based on a search query string

Flags
  - query: `--query "your youtube query here"`
  - maxResults: `--max-results 60`
  - requestsPerMin: `--requests-per-min 60`

## Search
Search from all the fetched youtube videos using an inverted index

Flags
  - query: `--query "you local query here"`

## Exploration Guide
`fetch.go` - source code for `fetch`

`search.go` - source code for `search`

`ii.go` - defines the InvertedIndex type among other types as well as relevant methods

`util.go` - contains the extractLowerAlphanumericFields definition
