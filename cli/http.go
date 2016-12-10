package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/goware/urlx"
)

func joinURL(base, path string) (string, error) {
	u, err := urlx.Parse(base)
	if err != nil {
		return "", err
	}

	splits := strings.Split(path, "?")

	u.Path = splits[0]
	if len(splits) > 1 {
		u.RawQuery = url.QueryEscape(strings.Join(splits[1:], ""))
	}

	return u.String(), nil
}

func httpRequest(method, url string, data interface{}, jwt string) (*http.Response, error) {
	const defaultTimeout = 10 * time.Second

	var contentType string
	var body io.Reader

	switch data.(type) {
	case nil:
		body = nil
		contentType = ""

	case string:
		contentType = "plain/text"
		body = strings.NewReader(data.(string))

	case []byte:
		contentType = "application/octet-stream"
		body = bytes.NewReader(data.([]byte))

	case io.Reader:
		contentType = "application/octet-stream"
		body = data.(io.Reader)

	default:
		contentType = "application/json"

		bodyRead, bodyWrite := io.Pipe()
		go func() {
			json.NewEncoder(bodyWrite).Encode(data)
			bodyWrite.Close()
		}()

		body = bodyRead
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if jwt != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))
	}

	if contentType != "" {
		req.Header.Add("content-type", contentType)
	}

	client := &http.Client{
		Timeout: defaultTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
