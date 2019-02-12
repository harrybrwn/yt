package youtube

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	// Logging is a variable serves as a toggle for builtin logging
	client = &http.Client{}
)

const (
	badchars = `\/:*?"<>|.`
	agent    = "Video download cli tool"
)

func get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s%d", agent, time.Now().Nanosecond()))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	return buf.Bytes(), err
}

func safeFileName(name string) string {
	for i := range badchars {
		if strings.Contains(name, string(badchars[i])) {
			name = strings.Replace(name, string(badchars[i]), "", -1)
		}
	}
	return name
}
