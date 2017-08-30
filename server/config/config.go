package config

type Config struct {
	Server struct {
		Addr       string `toml:"addr"`
		PublicAddr string `toml:"public_addr"`
		BundlesDir string `toml:"bundles_dir"`
	} `toml:"server"`

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
