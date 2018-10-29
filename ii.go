package main

import (
	"os"
	"strings"
	"io/ioutil"
	"path/filepath"
	"encoding/json"

	"google.golang.org/api/youtube/v3"
)

const videoDataDirPath = "./data/videos/"
const invertedIndexFile = "./data/ii.json"

// Keyword is a string query value
type Keyword string

// VideoID is an identifier for a youtube video
type VideoID string

// VideoIDs is a slice of VideoIDs
type VideoIDs []VideoID

// Thumbnail is a struct representing an image thumbnail
type Thumbnail struct {
	Height int
	URL string
	Width int
}

// Thumbnails is a struct representing a collection of thumbnails
type Thumbnails struct {
	Default Thumbnail
	High Thumbnail
	Medium Thumbnail
}

// Video represents the video data we collect from youtube
type Video struct {
	ChannelID string
	ChannelTitle string
	Description string
	LiveBroadcastContent string
	PublishedAt string
	Thumbnails Thumbnails
	Title string
}

// Exists checks if a VideoID is present in VideoIDs
func (vis VideoIDs) Exists(vi VideoID) bool {
	for _, val := range vis {
		if val == vi {
			return true
		}
	}
	return false
}

// InvertedIndex represents our particular use case of an inverted index - Keywords: VideoIds
type InvertedIndex map[Keyword]VideoIDs

// BuildFromFile path, which is a relative path to the inverted index data file
func (ii *InvertedIndex) BuildFromFile() error {
	file, err := os.Open(invertedIndexFile)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			return nil
		default:
			return err
		}
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(ii)
	if err != nil {
		return err
	}
	return nil
}

// Extend existing inverted index using a Youtube SearchResult
// write new video datum to file
func (ii *InvertedIndex) Extend(item *youtube.SearchResult) error {
	if _, err := os.Stat(videoDataDirPath); os.IsNotExist(err) {
	    os.MkdirAll(videoDataDirPath, os.ModePerm)
	}
	videoID := VideoID(item.Id.VideoId)
	titleDescription := strings.Join([]string{item.Snippet.Title, item.Snippet.Description}, " ")
	cleanTitleDescriptionFields := extractLowerAlphanumericFields(titleDescription)
	for _, fld := range cleanTitleDescriptionFields {
		keyword := Keyword(fld)
		videoIDs, arePresent := (*ii)[keyword]
		if arePresent && !videoIDs.Exists(videoID) {
			(*ii)[keyword] = append((*ii)[keyword], videoID)
			snippetB, err := json.MarshalIndent(item.Snippet, "", "  ")
			if err != nil {
				return err
			}
			videoDataFilename := strings.Join([]string{string(videoID), ".json"}, "")
			videoDataFilePath := filepath.Join(videoDataDirPath, videoDataFilename)
			err = ioutil.WriteFile(videoDataFilePath, snippetB, 0666)
			if err != nil {
				return err
			}
		} else {
			(*ii)[keyword] = VideoIDs{videoID}
		}
	}
	return nil
}

// WriteToDisk the inverted index
func (ii *InvertedIndex) WriteToDisk() error {
	iib, err := json.MarshalIndent(ii, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(invertedIndexFile, iib, 0666)
	if err != nil {
		return err
	}
	return nil
}

// Query the inverted index with keywords
func (ii *InvertedIndex) Query(query string) ([](*Video), error) {
	videoIDCtTracker := make(map[VideoID]int)
	cleanQueryFields := extractLowerAlphanumericFields(query)
	intersectionDepth := 0
	for _, field := range cleanQueryFields {
		videoIDs, arePresent := (*ii)[Keyword(field)]
		if arePresent {
			for _, videoID := range videoIDs {
				if _, isVideoIDInTracker := videoIDCtTracker[videoID]; isVideoIDInTracker {
					videoIDCtTracker[videoID]++
				} else {
					videoIDCtTracker[videoID] = 1
				}
			}
			intersectionDepth++
		}
	}
	videos := [](*Video){}
	for videoID, ct := range videoIDCtTracker {
		if ct == intersectionDepth {
			videoDataFilename := strings.Join([]string{string(videoID), ".json"}, "")
			videoDataFilePath := filepath.Join(videoDataDirPath, videoDataFilename)
			videoDataFile, err := os.Open(videoDataFilePath)
			if err != nil {
				return nil, err
			}
			defer videoDataFile.Close()
			video := &Video{}
			err = json.NewDecoder(videoDataFile).Decode(video)
			if err != nil {
				return nil, err
			}
			videos = append(videos, video)
		}
	}
	return videos, nil
}
