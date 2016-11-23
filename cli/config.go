package cli

import (
	"os"
	"runtime"

	"path"

	"github.com/BurntSushi/toml"
)

type config interface {
	Path() string
}

type projectConfig struct {
	// [server]
	Server struct {
		Addr string `toml:"addr"`
	} `toml:"server"`
}

func (p *projectConfig) Path() string {
	return "./WarpFile"
}

type globalConfig struct {
	Sessions [][]string `toml:"sessions"`
}

func (g *globalConfig) Path() string {
	var home string
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	return path.Join(home, ".warp")
}

func configLoad(conf config) error {
	_, err := toml.DecodeFile(conf.Path(), conf)
	return err
}

func configSave(conf config) error {
	file, err := os.Create(conf.Path())

	if err != nil {
		return err
	}

	defer file.Close()
	return toml.NewEncoder(file).Encode(conf)
}

func ProjectConfig() (*projectConfig, error) {
	conf := &projectConfig{}

	err := configLoad(conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func GlobalConfig() (*globalConfig, error) {
	conf := &globalConfig{}

	err := configLoad(conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
