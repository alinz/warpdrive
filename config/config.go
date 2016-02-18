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
