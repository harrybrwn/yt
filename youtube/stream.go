package youtube

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// Stream represents a meta-data container that also contains the link
// to the actual youtube video data.
type Stream struct {
	// to be honest I don't know what this is but it's important for
	// determining whether the stream is a video or audio stream, or both.
	MimeType MimeType
	// The url for the raw video data
	URL string `json:"url"`
	// Height of the video in pixels.
	// If Height is equal to zero, it is an audio-only stream.
	Height int `json:"height"`
	// Width of the video in pixels.
	// If Width is equal to zero, it is an audio-only stream.
	Width int `json:"width"`
	// Bitrate of the video's audio.
	Bitrate int `json:"bitrate"`
	// size of the video
	ContentLength string `json:"contentLength"`

	SignatureCipher string `json:"signatureCipher"`
}

// WriteTo will write the stream data to an io.Writer
func (s Stream) WriteTo(w io.Writer) (int64, error) {
	resp, err := client.Get(s.URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return io.Copy(w, resp.Body)
}

// IsDualStream returns true if the stream contains both audio and video
func (s Stream) IsDualStream() bool {
	return len(s.MimeType.Codecs) > 1
}

// IsAudioStream returns true if the stream contains audio
func (s Stream) IsAudioStream() bool {
	if len(s.MimeType.Codecs) > 1 {
		return true
	}
	return strings.Contains(s.MimeType.ContentType, "audio")
}

// IsVideoStream returns true if the stream contains video
func (s Stream) IsVideoStream() bool {
	if len(s.MimeType.Codecs) > 1 {
		return true
	}
	return strings.Contains(s.MimeType.ContentType, "video")
}

// GetURL will return are parsed version of the url. Returns nil on error.
func (s Stream) GetURL() (*url.URL, error) {
	var u = s.URL
	if s.URL == "" {
		q, err := url.ParseQuery(s.SignatureCipher)
		if err != nil {
			return nil, err
		}
		u = q.Get("url")
		if u == "" {
			return nil, errors.New("no url found")
		}
		return url.Parse(u)
	}
	return url.Parse(u)
}

// GetBestStream returns the stream with the hightest width and height.
// Only works for streams containing video.
func GetBestStream(c *[]Stream) *Stream {
	maxh, maxw := 0, 0
	var s Stream
	for _, strm := range *c {
		if strm.Height > maxh && strm.Width > maxw {
			maxh = strm.Height
			maxw = strm.Width
			s = strm
		}
	}
	return &s
}

// DownloadFromStream accepts a stream and downloads it to a given file name.
func DownloadFromStream(s *Stream, fname string) error {
	url, err := s.GetURL()
	if err != nil {
		return err
	}
	resp, err := client.Do(&http.Request{
		Method: "GET",
		Proto:  "HTTP/1.1",
		URL:    url,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

var (
	codecsRegex   = regexp.MustCompile(`;\s?codecs=`)
	mimeTypeRegex = regexp.MustCompile(`^"(.*?);\s?codecs=\\"(.*?)\\""$`)
)

// MimeType is a mimetype
type MimeType struct {
	ContentType string
	Codecs      []string
}

// UnmarshalJSON makes the MimeType struct implement the
// json.Unmarshaler interface
func (m *MimeType) UnmarshalJSON(b []byte) error {
	result := mimeTypeRegex.FindAllStringSubmatch(string(b), -1)
	m.ContentType = result[0][1]
	m.Codecs = strings.Split(result[0][2], ", ")
	return nil
}

var _ json.Unmarshaler = (*MimeType)(nil)

func splitMimeType(s string) (string, []string) {
	match := codecsRegex.Split(s, 2)
	return match[0], regexp.MustCompile(`,\s`).Split(match[1], 2)
}

// sorts streams into two groups, NOT by numeric value. The two groups
// are streams containing audio and streams containing video.
func sortStreams(streams []Stream) ([]Stream, []Stream) {
	var (
		vsc, asc int = 0, 0
		s        Stream
	)

	// find length of each []Stream
	for _, s = range streams {
		if s.IsVideoStream() {
			vsc++
		} else if s.IsAudioStream() {
			asc++
		}
	}

	vstreams := make([]Stream, vsc)
	astreams := make([]Stream, asc)
	vsc, asc = 0, 0
	for _, s = range streams {
		if s.IsVideoStream() {
			vstreams[vsc] = s
			vsc++
		} else if s.IsAudioStream() {
			astreams[asc] = s
			asc++
		}
	}
	return vstreams, astreams
}

var (
	_ io.WriterTo      = (*Stream)(nil)
	_ json.Unmarshaler = (*MimeType)(nil)
)
