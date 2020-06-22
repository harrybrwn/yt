package youtube

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) {
	// id := "FidhD-izZnk"
	// id := "O9Ks3_8Nq1s"
	// id := "9bZkp7q19f0"
}

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
	thumbnail := temp() + "_" + "thumbnail"
	if err = v.Thumbnails[0].Download(thumbnail); err != nil {
		t.Error(err)
	}
	os.Remove(thumbnail)
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
	err = &playabilityStatus{Status: "test", Reason: "testing"}
	if err.Error() != "test: testing" {
		t.Error("wrong error message")
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

func TestGetFromHTML(t *testing.T) {
	id := "9bZkp7q19f0"
	v, err := NewVideo(id)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range v.Streams {
		if _, err := s.GetURL(); err != nil {
			t.Error(err)
		}
	}
}

func TestSafeFileName(t *testing.T) {
	origin := `\/:*?"<>|.`
	safe := safeFileName(origin)
	if safe != "" {
		t.Errorf(`'safeFileName' is not getting rid of the right characters. Expected: "", got: "%s"`, safe)
	}
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
