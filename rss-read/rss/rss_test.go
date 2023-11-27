package rss_test

import (
	"os"
	"testing"
	"time"

	"github.com/hsmtkk/friendly-funicular/rss-read/rss"
	"github.com/stretchr/testify/assert"
)

const RSS_URL = "https://news.yahoo.co.jp/rss/topics/top-picks.xml"

func TestGetBytes(t *testing.T) {
	contents, err := rss.GetBytes(RSS_URL)
	assert.Nil(t, err)
	assert.True(t, len(contents) > 0)
}

func TestDecodeXML(t *testing.T) {
	contents, err := os.ReadFile("sample.xml")
	assert.Nil(t, err)
	got, err := rss.DecodeXML(contents)
	assert.Nil(t, err)
	assert.Equal(t, 8, len(got))
	assert.Equal(t, "北朝鮮 軍事衛星打ち上げを正当化", got[0].Title)
	assert.Equal(t, "https://news.yahoo.co.jp/pickup/6483056?source=rss", got[0].Link)
	assert.Greater(t, time.Now(), got[0].PubDate)
}
