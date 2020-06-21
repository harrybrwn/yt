package youtube

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	client = &http.Client{
		Transport: &userAgentTransport{
			agent: agent,
			inner: http.DefaultTransport,
		},
	}
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

type userAgentTransport struct {
	agent string
	inner http.RoundTripper
}

func (uat *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if uat.inner == nil {
		uat.inner = http.DefaultTransport
	}
	req.Header.Set("User-Agent", uat.agent)
	return uat.inner.RoundTrip(req)
}

func safeFileName(name string) string {
	for i := range badchars {
		if strings.Contains(name, string(badchars[i])) {
			name = strings.Replace(name, string(badchars[i]), "", -1)
		}
	}
	return name
}

type byteReader interface {
	ReadBytes(byte) ([]byte, error)
	ReadByte() (byte, error)
}

func parseQuery(r byteReader, m map[string][][]byte) (err error) {
	var (
		b          []byte
		key, val   string
		err1, err2 error
	)

	for {
		b, err = r.ReadBytes('&')
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return err
		}
		if i := bytes.Index(b, []byte{'='}); i >= 0 {
			key, val = string(b[:i]), string(b[i+1:len(b)-1])
		}
		if len(key) == 0 {
			continue
		}

		key, err1 = url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		val, err2 = url.QueryUnescape(string(val))
		if err2 != nil {
			if err == nil {
				err = err2
			}
			continue
		}
		m[key] = append(m[key], []byte(val))
	}
	return
}
