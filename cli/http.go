package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

func encodeJSON(data interface{}) io.Reader {
	bodyRead, bodyWrite := io.Pipe()

	go func() {
		json.NewEncoder(bodyWrite).Encode(data)
		bodyWrite.Close()
	}()

	return bodyRead
}

func decodeJSON(from io.Reader, to interface{}) error {
	return json.NewDecoder(from).Decode(to)
}

func httpRequest(method, url string, data interface{}) (*http.Response, error) {
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
		body = encodeJSON(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
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
