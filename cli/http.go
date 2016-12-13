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
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	baseU, err := urlx.Parse(base)
	if err != nil {
		return "", err
	}

	u.Scheme = baseU.Scheme
	u.Host = baseU.Host

	return u.String(), nil
}

func httpRequest(method, url string, data interface{}, jwt, contentType string) (*http.Response, error) {
	const defaultTimeout = 10 * time.Second
	var body io.Reader

	switch data.(type) {
	case nil:
		body = nil
		contentType = ""

	case string:
		if contentType == "" {
			contentType = "plain/text"
		}
		body = strings.NewReader(data.(string))

	case []byte:
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		body = bytes.NewReader(data.([]byte))

	case io.Reader:
		if contentType == "" {
			contentType = "multipart/octet-stream"
		}
		body = data.(io.Reader)

	default:
		if contentType == "" {
			contentType = "application/json"
		}

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
