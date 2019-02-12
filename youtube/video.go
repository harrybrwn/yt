package youtube

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		panic(err)
	}
	if !partialConfigREGEX.Match(b) {
		return nil, errors.New("couldn't find data")
	}
	rawb := bytes.Replace(
		bytes.Replace(
			bytes.Replace(
				partialConfigREGEX.FindAllSubmatch(b, 1)[0][1],
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
	if err != nil {
		return err
	}
	v.baseVideo = vd.VideoDetails.baseVideo
	v.Streams = vd.StreamingData.Formats
	vstream, astream := sortStreams(&vd.StreamingData.AdaptiveFormats)
	v.VideoStreams, v.AudioStreams = *vstream, *astream
	v.FileName = safeFileName(vd.VideoDetails.baseVideo.Title)
	return nil
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
	var vid Video
	data, err := getRaw(id)
	if err != nil {
		return nil, err
	}
	err = initVideoData(data, &vid)
	if err != nil {
		return nil, err
	}
	return &vid, nil
}

// Download will download the video given a file name.
//
// It is suggested that '.mp4' is used as the extention
// in the file name but is not manditory.
func (v *Video) Download(fname string) error {
	best := GetBestStream(&v.Streams)
	return DownloadFromStream(best, fname)
}

// DownloadAudio will download the video's audio given a file name.
//
// The suggested file extention is '.mpa'
func (v *Video) DownloadAudio(fname string) error {
	max := 0
	var high *Stream
	for _, s := range v.AudioStreams {
		if s.Bitrate > max {
			high = &s
			max = s.Bitrate
		}
	}
	return DownloadFromStream(high, fname)
}

// DownloadVideo will download the video given a file name.
//
// It is suggested that '.mp4' is used as the extention
// in the file name but is not manditory.
func (v *Video) DownloadVideo(fname string) error {
	max := 0
	var high *Stream
	for _, s := range v.VideoStreams {
		if s.Height > max {
			high = &s
			max = s.Height
		}
	}
	return DownloadFromStream(high, fname)
}
