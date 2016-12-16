package warpify

import (
	"path/filepath"
)

type warpConf struct {
	bundleVersion   string
	bundlePath      string
	documentPath    string
	productionName  string
	automaticUpdate bool
	pubSub          pubSub
}

func (c *warpConf) getDocumentPath(path string) string {
	return filepath.Join(c.documentPath, path)
}

func (c *warpConf) getBundlePath(path string) string {
	return filepath.Join(c.bundlePath, path)
}

// conf is a global warpConf for this package
var conf warpConf
