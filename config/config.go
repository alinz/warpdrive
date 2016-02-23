package config

import "github.com/BurntSushi/toml"

//Config configuration
type Config struct {
	//[server]
	Server struct {
		Bind string `toml:"bind"`
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
		FileMaxSize int64  `toml:"file_max_size"`
		TempFolder  string `toml:"temp_folder"`
	} `toml:"file_upload"`

	//[bundle]
	Bundle struct {
		BundlesFolder string `toml:"bundles_folder"`
	} `toml:"bundle"`

	//[static]
	Static struct {
		Path string `toml:"path"`
	} `toml:"static"`
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

	return config, nil
}
