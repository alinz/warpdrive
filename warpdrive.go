package warpdrive

import "github.com/pressly/warpdrive/config"

// The following two VERSION and LONGVERSION will be set during the build

var (
	// VERSION describes the version number - either dev or hash value
	VERSION string

	// LONGVERSION describes the version number - either dev or hash value
	LONGVERSION string

	// Conf Server Conf as global config variable.
	Conf *config.ServerConfig
)

func init() {
	if VERSION == "" {
		VERSION = "dev"
		LONGVERSION = "development"
	}
}
