package main

import (
	"log"
	"flag"
	"time"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const developerKey = "AIzaSyC6L_9R0Ji5ILgt09af93T0zNMolt69c5Q"
const numVideosToFetch = 1000

var (
	query = flag.String("query", "Google", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
	requestsPerMin = flag.Int64("requests-per-min", 60, "Limit request rate per minute")
	// 1 / req/min * sec/min * ms/sec * µs/ms * ns/µs -> ns/req
	rateLimitPeriod = time.Duration(1 / float64(*requestsPerMin) * 60 * 1000 * 1000 * 1000)
	numPagesToFetch = numVideosToFetch / int(*maxResults) // videos / videos/page -> pages
)

// ready youtube service
// read existing inverted index from file
// call service.Search.List while there are still pages to fetch
// extend inverted index for each call
// write inverted index to file
func main() {
	flag.Parse()
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	ii := InvertedIndex{}
	err = ii.BuildFromFile()
	if err != nil {
		log.Fatalf("Error building inverted index from file: %v", err)
	}
	var pageToken string
	rateLimiter := time.Tick(rateLimitPeriod)
	for numPagesToFetch > 0 {
		<-rateLimiter
		call := service.Search.List("id,snippet").
			Q(*query).
			MaxResults(*maxResults).
			PageToken(pageToken)
		log.Printf("delayed for %v, fetching page %v\n", rateLimitPeriod, numPagesToFetch)
		response, err := call.Do()
		if err != nil {
			log.Fatalf("Error calling YouTube api: %v", err)
		}
		pageToken = response.NextPageToken
		log.Printf("nextPageToken -> %v", pageToken)
		for _, item := range response.Items {
			if item.Id.Kind == "youtube#video" {
				err = ii.Extend(item)
				if err != nil {
					log.Fatalf("Error extending inverted index: %v", err)
				}
			}
		}
		numPagesToFetch--
	}
	err = ii.WriteToDisk()
	if err != nil {
		log.Fatalf("Error writing inverted index to disk: %v", err)
	}
}
