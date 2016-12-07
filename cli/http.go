package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

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
