package youtube

import (
	"testing"
)

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

	_, err = NewPlaylist("")
	if err == nil {
		t.Error("expected error")
	}
}
