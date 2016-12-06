package warp_test

import (
	"bytes"
	"fmt"
	"os"
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

	f, _ := os.Create("./output")
	defer f.Close()
	f.Write(buff.Bytes())

	f.Seek(0, 0)
	w = warp.NewReader(f)
	err = w.Extract("./ali")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(buff.Len())
}
