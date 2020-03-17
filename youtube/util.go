package youtube

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	client = &http.Client{}
)

const (
	badchars = `\/:*?"<>|.`
	agent    = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/500.0 (KHTML, like Gecko) Chrome/70.0.0.0 Safari/500.0"
)

func get(urlStr string) ([]byte, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: "GET",
		Host:   parsedURL.Host,
		Proto:  "HTTP/1.1",
		Header: http.Header{
			"User-Agent": []string{fmt.Sprintf("%s%d", agent, time.Now().Nanosecond())},
		},
		URL: parsedURL,
	}

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
