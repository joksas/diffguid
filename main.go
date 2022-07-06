package main

import (
	"os"

	"github.com/mmcdole/gofeed"
)

func localFeed(path string) (*gofeed.Feed, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return gofeed.NewParser().Parse(f)
}

func remoteFeed(url string) (*gofeed.Feed, error) {
	return gofeed.NewParser().ParseURL(url)
}
