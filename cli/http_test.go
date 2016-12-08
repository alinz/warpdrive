package cli

import "testing"

func TestJoinPath(t *testing.T) {
	testCases := []struct {
		base   string
		target string
		result string
	}{
		{"http://localhost:1234", "/apps", "http://localhost:1234/apps"},
		{"http://localhost:1234", "/apps/1", "http://localhost:1234/apps/1"},
		{"localhost:1234", "/apps", "http://localhost:1234/apps"},
		{"localhost:1234", "/apps/1", "http://localhost:1234/apps/1"},
		{"https://localhost:1234", "/apps?name=hello world", "https://localhost:1234/apps?name%3Dhello+world"},
	}

	for _, testCase := range testCases {
		path, err := joinURL(testCase.base, testCase.target)
		if err != nil {
			t.Error(err)
		}

		if testCase.result != path {
			t.Errorf("expected '%s', got '%s", testCase.result, path)
		}
	}

}
