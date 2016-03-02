package warp_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/pressly/warpdrive/lib/warp"
)

func TestWarpEncoding(t *testing.T) {
	var buff bytes.Buffer

	w := warp.NewWriter(&buff)

	err := w.AddFile("file1", "file1")
	if err != nil {
		t.Error(err)
	}

	err = w.AddFile("file2", "file2")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(buff.Bytes())
}
