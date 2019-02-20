package youtube

import (
	"testing"
)

func TestSafeFileName(t *testing.T) {
	origin := `\/:*?"<>|.`
	safe := safeFileName(origin)
	if safe != "" {
		t.Errorf(`'safeFileName' is not getting rid of the right characters. Expected: "", got: "%s"`, safe)
	}
}

func TestGet_Err(t *testing.T) {
	_, err := get("invalid")
	if err == nil {
		t.Error("expected error")
	}
}

func TestURLForatting(t *testing.T) {
	if videoURL("video id") != "http://www.youtube.com/watch?v=video id" {
		t.Error("invalid url format")
	}
	if playlistURL("another id") != "https://www.youtube.com/playlist?list=another id" {
		t.Error("invalid url format")
	}
}
