package entities

import (
	"fmt"
	"strconv"
)

type Media struct {
	ModelType   string
	ModelID     int
	FileName    string
	OrderColumn string
	MediaType   string
}

type AboutUsMedia struct {
	MediaType int    `json:"type"`
	URL       string `json:"url"`
}

func AboutUsFromMedia(media Media, host string) AboutUsMedia {
	mediaType, _ := strconv.Atoi(media.MediaType)
	return AboutUsMedia{
		MediaType: mediaType,
		URL:       media.getURL(host),
	}
}

func (m Media) getURL(host string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s", host, m.OrderColumn, m.FileName)
}
