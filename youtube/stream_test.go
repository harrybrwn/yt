package youtube

import (
	"testing"
)

func TestURLForatting(t *testing.T) {
	if videoURL("video id") != "http://www.youtube.com/watch?v=video id" {
		t.Error("invalid url format")
	}
	if playlistURL("another id") != "https://www.youtube.com/playlist?list=another id" {
		t.Error("invalid url format")
	}
}
