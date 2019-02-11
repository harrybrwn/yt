package youtube

import (
	"testing"
)

func TestNewVideo(t *testing.T) {
	_, err := NewVideo("_gvweJDAIGE")
	if err != nil {
		t.Error(err)
	}
}

func TestInitVideoData(t *testing.T) {
	// id := "jFXCindxT1M"
	// id := "_gvweJDAIGE"

	// f, err := os.Open("./testdata/video")
	// if err != nil {
	// 	t.Error(err)
	// }
	// defer f.Close()
	// b, err := ioutil.ReadAll(f)
	// if err != nil {
	// 	t.Error(err)
	// }

	// all := partialConfigREGEX.FindAllSubmatch(b, 1)
	// _ = getRaw(videoURL(id))
	// fmt.Println(string(raw))

	// rawb := bytes.Replace(
	// 	bytes.Replace(
	// 		bytes.Replace(
	// 			all[0][1],
	// 			[]byte("\\\\"),
	// 			[]byte(`ESCAPE`),
	// 			-1,
	// 		),
	// 		[]byte("\\"),
	// 		[]byte(""),
	// 		-1,
	// 	),
	// 	[]byte(`ESCAPE`),
	// 	[]byte("\\"),
	// 	-1,
	// )
	// raw := append([]byte("{"), append(rawb, "}"...)...)

	// v := Video{}
	// err := initVideoData(raw, &v)
	// if err != nil {
	// 	t.Error(err)
	// }
}
