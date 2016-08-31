package warpdrive

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

//Config the configuration of warpdrive app
type Config struct {
	//[server]
	Server struct {
		Addr          string `toml:"addr"`
		DataDir       string `toml:"data_dir"`
		BundlesFolder string
		TempFolder    string
	} `toml:"server"`

	//[jwt]
	JWT struct {
		SecretKey string `toml:"jwt_key"`
		MaxAge    int    `toml:"max_age"`
		Path      string `toml:"path"`
		Domain    string `toml:"domain"`
		Secure    bool   `toml:"secure"`
	} `toml:"jwt"`

	//[db]
	DB struct {
		Database string `toml:"database"`
		Hosts    string `toml:"hosts"`
		Username string `toml:"username"`
		Password string `toml:"password"`
	} `toml:"db"`

	//[security]
	Security struct {
		KeySize int `toml:"key_size"`
	} `toml:"security"`

	//[file_upload]
	FileUpload struct {
		FileMaxSize int64 `toml:"file_max_size"`
	} `toml:"file_upload"`
}

//Conf global config variable.
var Conf *Config

//NewConfig reads the config file and instanciated to global conf.
//you only need to call this one once. use it in your main app
func NewConfig(filename string) (*Config, error) {
	config := Config{}

	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return nil, err
	}

	if config.Server.DataDir == "" {
		return nil, errors.New("data_dir is not configured")
	}

	dd, err := os.Stat(config.Server.DataDir)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("data_dir: %s does not exist", config.Server.DataDir)
	}
	if !dd.IsDir() {
		return nil, fmt.Errorf("data_dir: %s is not a directory", config.Server.DataDir)
	}

	config.Server.BundlesFolder = filepath.Join(config.Server.DataDir, "bundles")
	dir, err := os.Stat(config.Server.BundlesFolder)
	if err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(config.Server.BundlesFolder, 0755); err != nil {
			return nil, err
		}
	} else if !dir.IsDir() {
		return nil, fmt.Errorf("data_dir: %s is not a directory", config.Server.BundlesFolder)
	}

	config.Server.TempFolder = filepath.Join(config.Server.DataDir, "tmp")
	dir, err = os.Stat(config.Server.TempFolder)
	if err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(config.Server.TempFolder, 0755); err != nil {
			return nil, err
		}
	} else if !dir.IsDir() {
		return nil, fmt.Errorf("data_dir: %s is not a directory", config.Server.TempFolder)
	}

	Conf = &config

	return Conf, nil
}