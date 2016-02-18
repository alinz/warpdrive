package warpdrive

import (
	"github.com/pressly/warpdrive/config"
	"upper.io/db.v2"
)

var (
	Version string
	Config  *config.Config
	DB      db.Database
)
