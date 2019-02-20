package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var (
	basePlisturl       = "https://www.youtube.com/playlist?list=%s"
	plistRendererRegex = regexp.MustCompile(`"playlistVideoListRenderer":({.*})`)
)

func playlistURL(id string) string {
	return fmt.Sprintf(basePlisturl, id)
}

func getPlaylistData(id string) ([]byte, error) {
	b, err := get(playlistURL(id))
	if err != nil {
		return nil, err
	}

	data := plistRendererRegex.FindAllSubmatch(b, 1)
	if data == nil {
		return nil, errors.New("could not find playlist data")
	}
	raw := data[0][1]

	opens, closes := 0, 0 // refering to number of open and closed curly braces
	for k := range raw {
		if raw[k] == '{' {
			opens++
		} else if raw[k] == '}' {
			closes++
		}
		if opens == closes {
			return raw[:k+1], nil
		}
	}
	return nil, errors.New("could not parse responce")
}

// Playlist represents a youtube playlist.
type Playlist struct {
	Contents []struct {
		PlaylistVideoRenderer struct {
			VideoID            string `json:"videoId"`
			LengthSeconds      string `json:"lengthSeconds"`
			NavigationEndpoint struct {
				WatchEndpoint struct {
					VideoID          string `json:"videoId"`
					PlaylistID       string `json:"playlistId"`
					Index            int    `json:"index"`
					StartTimeSeconds int    `json:"startTimeSeconds"`
				} `json:"watchEndpoint"`
			} `json:"navigationEndpoint"`
		} `json:"playlistVideoRenderer"`
	} `json:"contents"`
}

// NewPlaylist creates a playlist object from a playlist id.
func NewPlaylist(id string) (*Playlist, error) {
	var p Playlist
	data, err := getPlaylistData(id)
	if err != nil {
		return nil, err
	}
	return &p, json.Unmarshal(data, &p)
}

// VideoIds returns a channel containing all of the video ids in the playlist.
func (p *Playlist) VideoIds() chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for _, content := range p.Contents {
			c <- content.PlaylistVideoRenderer.VideoID
		}
	}()
	return c
}
