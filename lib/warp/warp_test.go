package warp_test

import (
	"os"
	"testing"

	"github.com/pressly/warpdrive/lib/warp"
)

func TestWarpEncodingDecoding(t *testing.T) {
	output, err := os.Create("./output")
	if err != nil {
		t.Errorf(err.Error())
	}

	defer output.Close()

	files := make(map[string]string)
	files["./file1"] = "./file1"
	files["./file2"] = "./file2"
	files["./file3"] = "./file3"

	err = warp.Compress(files, output)
	if err != nil {
		t.Errorf(err.Error())
	}

	input, err := os.Open("./output")
	if err != nil {
		t.Errorf(err.Error())
	}

	defer input.Close()

	err = warp.Extract(input, "./extract")
	if err != nil {
		t.Errorf(err.Error())
	}
}
