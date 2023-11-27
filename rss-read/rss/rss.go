package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Feed struct {
	Link    string
	Title   string
	PubDate time.Time
}

func GetFeeds(newsURL string) ([]Feed, error) {
	bs, err := GetBytes(newsURL)
	if err != nil {
		return nil, err
	}
	feeds, err := DecodeXML(bs)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func GetBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get failed: %w", err)
	}
	defer resp.Body.Close()
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}
	return contents, nil
}

type schema struct {
	RSS     xml.Name      `xml:"rss"`
	Channel schemaChannel `xml:"channel"`
}

type schemaChannel struct {
	Items []schemaItem `xml:"item"`
}

type schemaItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func DecodeXML(xmlBytes []byte) ([]Feed, error) {
	sch := schema{}
	if err := xml.Unmarshal(xmlBytes, &sch); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal failed: %w", err)
	}
	results := []Feed{}
	for _, i := range sch.Channel.Items {
		parsedPubDate, err := time.Parse(time.RFC1123, i.PubDate)
		if err != nil {
			fmt.Printf("time.Parse failed: %s: %v\n", i.PubDate, err)
		}
		results = append(results, Feed{Link: i.Link, Title: i.Title, PubDate: parsedPubDate})
	}
	return results, nil
}
