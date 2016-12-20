package warpify

import (
	"path/filepath"

	"github.com/pressly/warpdrive/config"
)

type warpConf struct {
	bundleVersion string
	bundlePath    string
	documentPath  string
	platform      string
	defaultCycle  string
	forceUpdate   bool
	reloadTask    Callback
	pubSub        pubSub
	api           *api
	warpFile      *config.ClientConfig
}

func (c *warpConf) getDocumentPath(path string) string {
	return filepath.Join(c.documentPath, path)
}

func (c *warpConf) getBundlePath(path string) string {
	return filepath.Join(c.bundlePath, path)
}

// conf is a global warpConf for this package
var conf warpConf
