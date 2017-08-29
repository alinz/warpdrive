package config

type Config struct {
	Addr string `toml:"addr"`

	DB struct {
		Path string `toml:"path"`
	} `toml:"db"`

	TLS struct {
		CA      string `toml:"ca"`
		Private string `toml:"private"`
		Public  string `toml:"string"`
	} `toml:"tls"`

	Admin struct {
		Username string `toml:"username"`
		Password string `toml:"password"`
	} `toml:"admin"`
}
