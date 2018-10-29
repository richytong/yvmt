package main

import (
	"log"
	"flag"
	"encoding/json"
)

var (
	query = flag.String("query", "Google", "Search term")
)

func main() {
	flag.Parse()
	ii := InvertedIndex{}
	err := ii.BuildFromFile()
	if err != nil {
		log.Fatalf("Error building inverted index from file: %v", err)
	}
	videos, err := ii.Query(*query)
	if err != nil {
		log.Fatalf("Error querying inverted index: %v", err)
	}
	videosB, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println(string(videosB))
}
