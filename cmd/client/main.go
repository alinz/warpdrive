package warpdrive

import "fmt"
import "path"

type Callback interface {
	Do(value string, err error)
}

var (
	targetBundlePath   string
	targetDocumentPath string
	targetPlatform     string
	targetCallback     Callback
)

// Setup sets couple of variables which will be used internally
// to save and load bundles
func Setup(bundle, document, platform string) {
	targetBundlePath = bundle
	targetDocumentPath = document
	targetPlatform = platform

	fmt.Println("BundlePath: ", targetBundlePath)
	fmt.Println("DocumentPath", targetDocumentPath)

	warpFile := path.Join(targetBundlePath, "WarpFile")
	loadFile(warpFile)

	tempFile := path.Join(targetDocumentPath, "haha")
	//saveFile(tempFile, "Hello this is haha file")

	loadFile(tempFile)
}

// SourceBundlePath loads `version.warp` to check the version,
// - if the file does exists in `targetDocumentPath` then we know that we need to load the bundles from `targetDocumentPath/{version}`
// - if the file does not exists in `targetDocumentPath` then we simply load the file which bundle with app.
func SourceBundlePath() string {
	return ""
}

// Check starts the proces the process consists of couple of tasks.
// - we need to load the `WarpFile`, this file contains the basic information about where we can talk to Warpdrive server
// - then we need to check the current version
// 	-- check if a `version.warp` exists in `targetDocumentPath`, if not check the file in `targetBundlePath`
// - we make a GET request to server with version number by calling `/apps/x/cycles/x/releases/latest/version/x/platform/x`
// - we compare the version with current version
//  -- if they are equal, do nothing and return error saying nothing to update
// - generate a session key and encrypted with given `public key` in `WarpFile`
// - get the releaseId from response and POST it to server with the following api `/apps/x/cycles/x/releases/x/download`
// - decrypt the response using the previously generated session key and decode Warp file into `targetDocumentPath/{version}`
// - send a signal which restart the whole process.
func Check() {

}

func RegisterCallback(callback Callback) {
	targetCallback = callback
}
