package youtube

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// Playlist is a youtube playlist
type Playlist struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Views       int    `json:"views"`
	Description string `json:"description"`
	Videos      []struct {
		Author                       string  `json:"author"`
		Privacy                      string  `json:"privacy"`
		Comments                     string  `json:"comments"`
		Keywords                     string  `json:"keywords"`
		TimeCreated                  int     `json:"time_created"`
		Rating                       float64 `json:"rating"`
		Added                        string  `json:"added"`
		Likes                        int     `json:"likes"`
		CcLicense                    bool    `json:"cc_license"`
		CategoryID                   int     `json:"category_id"`
		SessionData                  string  `json:"session_data"`
		IsHd                         bool    `json:"is_hd"`
		EndscreenAutoplaySessionData string  `json:"endscreen_autoplay_session_data"`
		UserID                       string  `json:"user_id"`
		Title                        string  `json:"title"`
		Views                        string  `json:"views"`
		LengthSeconds                int     `json:"length_seconds"`
		Thumbnail                    string  `json:"thumbnail"`
		IsCc                         bool    `json:"is_cc"`
		Duration                     string  `json:"duration"`
		ID                           string  `json:"encrypted_id"`
		Description                  string  `json:"description"`
		Dislikes                     int     `json:"dislikes"`
	} `json:"video"`
}

// NewPlaylist creates a playlist object from a playlist id.
func NewPlaylist(id string) (*Playlist, error) {
	var req = http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       host,
		URL: &url.URL{
			Scheme: "https",
			Host:   host,
			Path:   "/list_ajax",
			RawQuery: url.Values{
				"list":            {id},
				"action_get_list": {"1"},
				"style":           {"json"},
				"hl":              {"en"},
				"index":           {"0"},
			}.Encode(),
		},
	}
	resp, err := client.Do(&req)
	if err != nil {
		return nil, err
	}
	var e error
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		var p Playlist
		return &p, json.NewDecoder(resp.Body).Decode(&p)
	default:
		e = &respError{}
	}
	if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return nil, err
	}
	return nil, e
}

type respError struct {
	Errors []string `json:"errors"`
}

func (e *respError) Error() string {
	return strings.Join(e.Errors, ", ")
}

// VideoIds returns a channel containing all of the video ids in the playlist.
func (p *Playlist) VideoIds() chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for _, vid := range p.Videos {
			c <- vid.ID
		}
	}()
	return c
}
