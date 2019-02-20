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

func TestGetRaw(t *testing.T) {
	b, err := getRaw("ferZnZ0_rSM")
	if err != nil {
		t.Error(err)
	}
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
