package youtube

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var (
	fullConfigREGEX    = regexp.MustCompile(`;ytplayer\.config\s*=\s*({.*?});`)
	partialConfigREGEX = regexp.MustCompile(`"player_response":"{(.*)}"`)
)

// Video represents a youbube video.
type Video struct {
	baseVideo

	// FileName is a file system safe version of the video's title.
	FileName string

	// A slice of stream objects containing both audio and video
	Streams []Stream

	// A slice of streams containing only video
	VideoStreams []Stream

	// A slice of streams containing only audio
	AudioStreams []Stream

	downloadStream Stream
}

// NewVideo creates and returns a new Video object.
func NewVideo(id string) (*Video, error) {
	vid := &Video{}
	r, err := info(id)
	if err != nil {
		return nil, err
	}
	defer r.cleanup()

	query := make(map[string][][]byte)

	if err = parseQuery(r, query); err != nil {
		return nil, err
	}
	if errcode := query["errorcode"]; len(errcode) != 0 {
		return nil, fmt.Errorf("could not find video %s: %s", id, errcode[0])
	}

	pResp, ok := query["player_response"]
	if !ok || len(pResp) < 1 {
		return nil, errors.New("could not find video player data")
	}

	return vid, initVideoData(pResp[0], vid)
}

// Download will download the video given a file name.
//
// It is suggested that '.mp4' is used as the extension
// in the file name but is not manditory.
func (v *Video) Download(fname string) error {
	s := GetBestStream(&v.Streams)
	return DownloadFromStream(s, fname)
}

// DownloadAudio will download the video's audio given a file name.
//
// The suggested file extension is '.mpa'
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
	defer r.cleanup()

	query := make(map[string][][]byte)
	return query, parseQuery(r, query)
}

type inforeader struct {
	*bufio.Reader
	cleanup func() error
}

func info(id string) (*inforeader, error) {
	req := &http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       "www.youtube.com",
		URL: &url.URL{
			Scheme:   "https",
			Host:     "www.youtube.com",
			Path:     "/get_video_info",
			RawQuery: url.Values{"video_id": {id}}.Encode(),
		},
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &inforeader{
		bufio.NewReader(resp.Body),
		func() error {
			return resp.Body.Close()
		},
	}, nil

}

func initVideoData(in []byte, v *Video) error {
	vd := VideoData{}
	err := json.Unmarshal(in, &vd)

	if vd.PlayabilityStatus.Status != "OK" {
		return fmt.Errorf("%s: %s",
			vd.PlayabilityStatus.Status, vd.PlayabilityStatus.Reason)
	}

	v.baseVideo = vd.VideoDetails.baseVideo
	v.Streams = vd.StreamingData.Formats
	vstream, astream := sortStreams(&vd.StreamingData.AdaptiveFormats)
	v.VideoStreams, v.AudioStreams = *vstream, *astream
	v.FileName = safeFileName(vd.VideoDetails.baseVideo.Title)
	return err
}
