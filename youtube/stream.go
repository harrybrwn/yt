package youtube

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// Stream represents a meta-data container that also contains the link
// to the actual youtube video data.
type Stream struct {
	// to be honest I don't know what this is but it's important for
	// deturmining whether the stream is a video or audio stream, or both.
	MimeType string `json:"mimeType"`

	// The url for the raw video data
	URL string `json:"probeUrl"`
	// URL string `json:"url"`

	// Height of the video in pixels.
	// If Height is equeal to zero, it is an audio-only stream.
	Height int `json:"height"`

	// Width of the video in pixels.
	// If Width is equeal to zero, it is an audio-only stream.
	Width int `json:"width"`

	// Bitrate of the video's audio.
	Bitrate int `json:"bitrate"`

	// size of the video
	ContentLength string `json:"contentLength"`
}

// IsDualStream returns true if the stream contains both audio and video
func (s Stream) IsDualStream() bool {
	_, codecs := splitMimeType(s.MimeType)
	return len(codecs) > 1
}

// IsAudioStream returns true if the stream contains audio
func (s Stream) IsAudioStream() bool {
	if s.IsDualStream() {
		return true
	}
	contentType, _ := splitMimeType(s.MimeType)
	return strings.Contains(contentType, "audio")
}

// IsVideoStream returns true if the stream contains video
func (s Stream) IsVideoStream() bool {
	if s.IsDualStream() {
		return true
	}
	contentType, _ := splitMimeType(s.MimeType)
	return strings.Contains(contentType, "video")
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
	b, err := get(s.URL)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fname, b, 0644)
}

func splitMimeType(s string) (string, []string) {
	match := regexp.MustCompile(`;\scodecs=`).Split(s, 2)
	return match[0], regexp.MustCompile(`,\s`).Split(match[1], 2)
}

// sorts streams into two groups, NOT by some numeric value. The two groups
// are streams containing audio and streams containing video.
func sortStreams(streams *[]Stream) (*[]Stream, *[]Stream) {
	var (
		vsc, asc int = 0, 0
		s        Stream
	)

	// find length of each []Stream
	for _, s = range *streams {
		if s.IsVideoStream() {
			vsc++
		} else if s.IsAudioStream() {
			asc++
		}
	}

	vstreams := make([]Stream, vsc)
	astreams := make([]Stream, asc)
	vsc, asc = 0, 0
	for _, s = range *streams {
		if s.IsVideoStream() {
			vstreams[vsc] = s
			vsc++
		} else if s.IsAudioStream() {
			astreams[asc] = s
			asc++
		}
	}
	return &vstreams, &astreams
}
