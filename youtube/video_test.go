package youtube

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewVideo(t *testing.T) {
	v, err := NewVideo("Nq5LMGtBmis")
	if err != nil {
		t.Error(err)
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
}

func TestDownloads(t *testing.T) {
	v, err := NewVideo("O9Ks3_8Nq1s")
	if err != nil {
		t.Error(err)
	}
	err = v.Download(temp())
	if err != nil {
		t.Error(err)
	}
	t.Run("audio download", func(t *testing.T) {
		err = v.DownloadAudio(temp())
		if err != nil {
			t.Error(err)
		}
		fmt.Println("ending audio download")
	})
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

// func testInitVideo(t *testing.T) {
// 	id := "Nq5LMGtBmis"
// 	b, _ := get(videoURL(id))
// 	ioutil.WriteFile("testingdata.json", partialConfigREGEX.FindAllSubmatch(b, 1)[0][1], 777)
// 	// vd := VideoData{}
// 	vd := map[string]interface{}{}

// 	data, err := getRaw(id)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if err := json.Unmarshal(data, &vd); err != nil {
// 		t.Error(err)
// 	}
// 	printJSON(vd)
// }

func TestGetRaw(t *testing.T) {
	b, err := getRaw("3eN7LKFCI8c")
	if err != nil {
		t.Error(err)
	}
	ioutil.WriteFile("testingdata.json", b, 777)
	if b == nil {
		t.Error("got empty byte array")
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
