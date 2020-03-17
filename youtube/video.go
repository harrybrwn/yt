package youtube

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

const baseURL = "http://www.youtube.com/watch?v=%s"

func videoURL(id string) string {
	return fmt.Sprintf(baseURL, id)
}

var (
	fullConfigREGEX    = regexp.MustCompile(`;ytplayer\.config\s*=\s*({.*?});`)
	partialConfigREGEX = regexp.MustCompile(`"player_response":"{(.*)}"`)
)

// getRaw is the function that retrieves the data from a youtube video
// and returns it as a channel containing a raw byte slice.
func getRaw(id string) ([]byte, error) {
	b, err := get(videoURL(id))
	if err != nil {
		return nil, err
	}
	if !partialConfigREGEX.Match(b) {
		return nil, errors.New("couldn't find data")
	}
	rawb := bytes.Replace(
		bytes.Replace(
			bytes.Replace(
				// partialConfigREGEX.FindAllSubmatch(b, 1)[0][1],
				fullConfigREGEX.FindAllSubmatch(b, 1)[0][1],
				[]byte("\\\\"),
				[]byte(`ESCAPE`),
				-1,
			),
			[]byte("\\"),
			[]byte(""),
			-1,
		),
		[]byte(`ESCAPE`),
		[]byte("\\"),
		-1,
	)
	return append([]byte("{"), append(rawb, "}"...)...), nil
}

func initVideoData(in []byte, v *Video) error {
	vd := VideoData{}
	err := json.Unmarshal(in, &vd)

	v.baseVideo = vd.VideoDetails.baseVideo
	v.Streams = vd.StreamingData.Formats
	vstream, astream := sortStreams(&vd.StreamingData.AdaptiveFormats)
	v.VideoStreams, v.AudioStreams = *vstream, *astream
	v.FileName = safeFileName(vd.VideoDetails.baseVideo.Title)
	return err
}

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
	info, err := info(id)
	if err != nil {
		return nil, err
	}
	vals, err := url.ParseQuery(info.String())
	if err != nil {
		return nil, err
	}

	pResp, ok := vals["player_response"]
	if !ok || len(pResp) < 1 {
		return nil, errors.New("could not find video player data")
	}

	return vid, initVideoData([]byte(pResp[0]), vid)
}

// Download will download the video given a file name.
//
// It is suggested that '.mp4' is used as the extention
// in the file name but is not manditory.
func (v *Video) Download(fname string) error {
	s := GetBestStream(&v.Streams)
	return DownloadFromStream(s, fname)
}

// DownloadAudio will download the video's audio given a file name.
//
// The suggested file extention is '.mpa'
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

func info(id string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

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
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
