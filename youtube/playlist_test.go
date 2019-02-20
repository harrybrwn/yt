package youtube

import (
	"testing"
)

func TestGetPlaylistData(t *testing.T) {
	id := "PLFsQleAWXsj_4yDeebiIADdH5FMayBiJo"
	b, err := getPlaylistData(id)
	if err != nil {
		t.Error(err)
	}
	if b == nil {
		t.Error("empty bytes")
	}
}

func TestNewPlaylist(t *testing.T) {
	id := "PLFsQleAWXsj_4yDeebiIADdH5FMayBiJo"
	p, err := NewPlaylist(id)
	if err != nil {
		t.Error(err)
	}
	if p == nil {
		t.Error("got <nil> playlist")
	}

	c := 0
	for ID := range p.VideoIds() {
		c++
		_ = ID
	}
	if c != len(p.Contents) {
		t.Error("different number of videos")
	}

	_, err = NewPlaylist("")
	if err == nil {
		t.Error("expected error")
	}
}
