package warpify

type config struct {
	bundleVersion string
	bundlePath    string
	documentPath  string
	pubSub        pubSub
}

var conf config
