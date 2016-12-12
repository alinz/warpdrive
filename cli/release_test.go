package cli

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBundleReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Create("./.test/test1")
		if err != nil {
			panic(err)
		}

		defer file.Close()

		io.Copy(file, r.Body)
		fmt.Println("just copied all the content")
	}))
	defer ts.Close()

	r, err := bundleReader("ios")
	if err != nil {
		panic(err)
	}

	_, err = httpRequest("POST", ts.URL, r, "", "")
	if err != nil {
		panic(err)
	}

}
