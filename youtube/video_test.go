package youtube

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewVideo(t *testing.T) {
	v, err := NewVideo("Nq5LMGtBmis")
	if err != nil {
		fmt.Printf("%T\n", err)
		t.Fatal(err)
	}
	if v == nil {
		t.Fatal("got nil video")
	}
	if v.ID != "Nq5LMGtBmis" {
		t.Error("wrong id")
	}
	if v.Author != "Vulf" {
		t.Error("wrong author")
	}
	if len(v.Streams) == 0 {
		t.Error("no streams")
	}
	if len(v.VideoStreams) == 0 {
		t.Error("no video steams")
	}
	if len(v.AudioStreams) == 0 {
		t.Error("no audio streams")
	}
	if v.ChannelID == "" {
		t.Error("should have non empty channel id")
	}
	for _, s := range v.Streams {
		if s.URL == "" {
			t.Error("stream has empty URL")
		}
	}
}

func TestMimeType(t *testing.T) {
	// v, err := NewVideo("O9Ks3_8Nq1s")
	v, err := NewVideo("FBkZ2TJZZUY")
	if err != nil {
		t.Fatal(err)
	}
	streams := make([]Stream, 0)
	streams = append(streams, v.Streams...)
	// streams = append(streams, v.AudioStreams...)
	// streams = append(streams, v.VideoStreams...)
	for _, s := range streams {
		fmt.Println(s.MimeType)
		// fmt.Printf("`%s`,\n", s.MimeType)
		// fmt.Println(splitMimeType(s.MimeType))
	}
}

func TestDownloads(t *testing.T) {
	v, err := NewVideo("O9Ks3_8Nq1s")
	if err != nil {
		t.Error(err)
	}
	if v == nil {
		t.Fatal("video should not be nil")
	}
	file := temp() + "_" + v.FileName
	err = v.Download(file)
	if err != nil {
		t.Error(err)
	}
	if err = os.Remove(file); err != nil {
		t.Error(err)
	}
	file = temp() + "_" + v.FileName
	t.Run("audio download", func(t *testing.T) {
		err = v.DownloadAudio(file)
		if err != nil {
			t.Error(err)
		}
	})
	if err = os.Remove(file); err != nil {
		t.Error(err)
	}
}

func TestVideo_Err(t *testing.T) {
	_, err := NewVideo("")
	if err == nil {
		t.Error("expected error")
	}
	err = initVideoData([]byte(""), &Video{})
	if err == nil {
		t.Error("expected error")
	}
	_, err = NewVideo("notavalidID")
	if err == nil {
		t.Error("expected error")
	}
}

func TestInfo(t *testing.T) {
	r, err := info("O9Ks3_8Nq1s")
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	m := make(map[string][][]byte)
	err = parseQuery(r, m)
	if err != nil {
		t.Error(err)
	}
	if len(m) == 0 {
		t.Error("should not have zero length")
	}
}

func Test(t *testing.T) {
	// id := "FidhD-izZnk"
	id := "O9Ks3_8Nq1s"
	// NewVideo(id)
	v := &Video{}
	r, err := info(id)
	if err != nil {
		t.Error(err)
	}

	query := make(map[string][][]byte)
	parseQuery(r, query)
	in := query["player_response"][0]
	if err = initVideoData(in, v); err != nil {
		t.Error(err)
	}
	// if err = v.Download("./test.mp4"); err != nil {
	// 	t.Error(err)
	// }
	// v.Thumbnails[0].Download("thumbnail")
	s := GetBestStream(&v.Streams)
	var b bytes.Buffer
	if _, err = s.WriteTo(&b); err != nil {
		t.Error(err)
	}
}

func TestMp4Tags(t *testing.T) {
}

func printJSON(m map[string]interface{}) {
	for key := range m {
		switch m[key].(type) {
		case map[string]interface{}:
			printJSON(m[key].(map[string]interface{}))
		default:
			fmt.Println(key, ": ", m[key])
		}
	}
}

func temp() string {
	f, err := ioutil.TempFile("", "yt")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
