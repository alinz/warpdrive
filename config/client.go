package config

import (
	"fmt"
)

// ClientConfig is being used in both cli and client bundle in mobile app
type ClientConfig struct {
	ServerAddr string `json:"server_addr"`
	AppID      int64  `json:"app_id"`
	AppName    string `json:"app_name"`
	CycleID    int64  `json:"cycle_id"`
	CycleName  string `json:"cycle_name"`
	Key        string `json:"key"`
}

// ClientConfigs collection of client's config
type ClientConfigs struct {
	configList []*ClientConfig
}

// Add new config into the list, make sure there are no duplicates
func (c *ClientConfigs) Add(newConfig *ClientConfig) error {
	for _, config := range c.configList {
		if config.AppID == newConfig.AppID &&
			config.CycleID == newConfig.CycleID &&
			config.ServerAddr == newConfig.ServerAddr {
			return fmt.Errorf("Configure app '%s' and cycle '%s' duplicated", config.AppName, config.CycleName)
		}
	}

	c.configList = append(c.configList, newConfig)

	return nil
}

// Load the config from file name with the given path
func (c *ClientConfigs) Load(path string) error {
	return nil
}

// Save config into a file in given path
func (c *ClientConfigs) Save(path string) error {
	return nil
}
