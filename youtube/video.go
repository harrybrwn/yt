package youtube

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

var (
	// ErrNoPlayerData is an error returned if the player data for a video
	// could not be found
	ErrNoPlayerData = errors.New("could not find player data")

	host               = "www.youtube.com"
	fullConfigREGEX    = regexp.MustCompile(`;ytplayer\.config\s*=\s*({.*?});`)
	partialConfigREGEX = regexp.MustCompile(`"player_response":"{(.*)}"`)
)

// Video id a youtube video.
type Video struct {
	baseVideo

	// FileName is a file system safe version of the video's title.
	FileName string
	// A slice of stream objects containing both audio and video
	Streams []Stream
	// A slice of streams containing only video
	VideoStreams []Stream
	// A slice of streams containing only audio
	AudioStreams   []Stream
	Thumbnails     []Thumbnail
	downloadStream Stream
}

// NewVideo creates and returns a new Video object.
func NewVideo(id string) (*Video, error) {
	vid := &Video{}
	r, err := info(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	query := make(map[string][][]byte)
	if err = parseQuery(r, query); err != nil {
		return nil, err
	}
	if errcode := query["errorcode"]; len(errcode) != 0 {
		return nil, fmt.Errorf("could not find video %s: %s", id, errcode[0])
	}

	pResp, ok := query["player_response"]
	if !ok || len(pResp) < 1 {
		return nil, ErrNoPlayerData
	}
	return vid, initVideoData(pResp[0], vid)
}

// Download will download the video given a file name.
//
// It is suggested that '.mp4' is used as the extension
// in the file name but is not mandatory.
func (v *Video) Download(fname string) error {
	s := GetBestStream(&v.Streams)
	return DownloadFromStream(s, fname)
}

// DownloadAudio will download the video's audio given a file name.
func (v *Video) DownloadAudio(fname string) error {
	var (
		max  = 0
		high *Stream
	)

	if len(v.AudioStreams) == 0 {
		return errors.New("no audio streams")
	}
	for _, s := range v.AudioStreams {
		if s.Bitrate > max {
			high = &s
			max = s.Bitrate
		}
	}
	return DownloadFromStream(high, fname)
}

// GetInfo returns a map of low-level video information used by youtube.
func GetInfo(id string) (map[string][][]byte, error) {
	r, err := info(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	query := make(map[string][][]byte)
	return query, parseQuery(r, query)
}

// Thumbnail is a video thumbnail
type Thumbnail struct {
	Height int
	Width  int
	URL    string
}

// Download will download the thumbnail to a file on disk
func (t *Thumbnail) Download(filename string) error {
	resp, err := client.Get(t.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

type inforeader struct {
	*bufio.Reader
	cleanup func() error
}

type byteReaderCloser interface {
	byteReader
	io.ReadCloser
}

func (ir *inforeader) Close() error {
	return ir.cleanup()
}

func info(id string) (byteReaderCloser, error) {
	var req = http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       host,
		URL: &url.URL{
			Scheme:   "https",
			Host:     host,
			Path:     "/get_video_info",
			RawQuery: url.Values{"video_id": {id}}.Encode(),
		},
	}
	resp, err := client.Do(&req)
	if err != nil {
		return nil, err
	}
	return &inforeader{
		bufio.NewReader(resp.Body),
		resp.Body.Close,
	}, nil
}

func initVideoData(in []byte, v *Video) (err error) {
	vd := VideoData{}
	if err = json.Unmarshal(in, &vd); err != nil {
		return err
	}
	v.baseVideo = vd.VideoDetails.baseVideo
	v.Streams = vd.StreamingData.Formats
	v.VideoStreams, v.AudioStreams = sortStreams(vd.StreamingData.AdaptiveFormats)
	v.FileName = safeFileName(vd.VideoDetails.baseVideo.Title)
	v.Thumbnails = vd.VideoDetails.Thumbnail.Thumbnails
	if vd.PlayabilityStatus.Status != "OK" {
		err = vd.PlayabilityStatus
	}
	return err
}
