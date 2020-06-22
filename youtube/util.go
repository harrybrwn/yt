package youtube

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
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

type userAgentTransport struct {
	agent string
	inner http.RoundTripper
}

func (uat *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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
		b        []byte
		key, val string
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

		key, err = url.QueryUnescape(key)
		if err != nil {
			return err
		}
		val, err = url.QueryUnescape(string(val))
		if err != nil {
			return err
		}
		m[key] = append(m[key], []byte(val))
	}
	return
}
