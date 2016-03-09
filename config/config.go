package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

//Config configuration
type Config struct {
	//[server]
	Server struct {
		Bind          string `toml:"bind"`
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
		Database string   `toml:"database"`
		Hosts    []string `toml:"hosts"`
		Username string   `toml:"username"`
		Password string   `toml:"password"`
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

//Load New read a configuration file and returns a Config object
func Load(configFile string, confEnv string) (*Config, error) {
	config := &Config{}

	if configFile == "" {
		configFile = confEnv
	}

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}

	if config.Server.DataDir == "" {
		return nil, errors.New("data_dir is not configured.")
	}

	dd, err := os.Stat(config.Server.DataDir)
	if err != nil && os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("data_dir: %s does not exist.", config.Server.DataDir))
	}
	if !dd.IsDir() {
		return nil, errors.New(fmt.Sprintf("data_dir: %s is not a directory.", config.Server.DataDir))
	}

	config.Server.BundlesFolder = filepath.Join(config.Server.DataDir, "bundles")
	if err := os.Mkdir(config.Server.BundlesFolder, 0755); err != nil {
		return nil, err
	}
	config.Server.TempFolder = filepath.Join(config.Server.DataDir, "tmp")
	if err := os.Mkdir(config.Server.TempFolder, 0755); err != nil {
		return nil, err
	}

	return config, nil
}
