package warpify

import (
	"testing"

	"github.com/pressly/warpdrive/lib/warp"
)

const testBundlePath = "./test/bundleFolder"
const testDocumentPath = "./test/documentFolder"

func TestDownloadRelease(t *testing.T) {
	Setup("1.0.0-prod", testBundlePath, testDocumentPath, "ios", "prod", true)

	r, err := conf.api.downloadVersion(1, 1, 1)
	if err != nil {
		panic(err)
	}

	warp.Extract(r, testDocumentPath)
}
