package config

// Config is base struct for reading configuration
type Config struct {
	DBPath     string `toml:"db_path"`
	BundlePath string `toml:"bundle_path"`

	Command struct {
		CAPath   string `toml:"ca_path"`
		CertPath string `toml:"cert_path"`
		KeyPath  string `toml:"key_path"`
		Addr     string `toml:"addr"`
	} `toml:"command"`

	Query struct {
		CAPath   string `toml:"ca_path"`
		CertPath string `toml:"cert_path"`
		KeyPath  string `toml:"key_path"`
		Addr     string `toml:"addr"`
	} `toml:"query"`
}
