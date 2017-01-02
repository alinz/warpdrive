package warpify

import "github.com/pressly/warpdrive/config"

type warpConf struct {
	bundleVersion string
	bundlePath    string
	documentPath  string
	platform      string
	defaultCycle  string
	forceUpdate   bool
	reloadTask    Callback
	api           *api
	warpFile      *config.ClientConfig
	reloadReady   chan struct{}
	chanClosed    bool
}

// conf is a global warpConf for this package
var conf *warpConf

func init() {
	// we need to initialize reloadReady before the application started
	conf = &warpConf{
		reloadReady: make(chan struct{}),
		chanClosed:  false,
	}
}
