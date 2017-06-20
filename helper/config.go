package helper

import (
	"os"

	"github.com/BurntSushi/toml"
)

// ConfigFile reads the config file
func ConfigFile(path string, conf interface{}) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	_, err = toml.DecodeFile(path, conf)
	if err != nil {
		return err
	}

	return nil
}
