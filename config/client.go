package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// These two path variables are being used for cli only
const iosWarpFilePath = "ios/WarpFile"
const androidWarpFilePath = "android/app/src/main/assets/WarpFile"

// ClientConfig is being used in both cli and client bundle in mobile app
type ClientConfig struct {
	ServerAddr string `json:"server_addr"`
	AppID      int64  `json:"app_id"`
	AppName    string `json:"app_name"`
	CycleID    int64  `json:"cycle_id"`
	CycleName  string `json:"cycle_name"`
	Key        string `json:"key"`
}

func (c *ClientConfig) equal(a *ClientConfig) bool {
	return c.AppID == a.AppID &&
		c.AppName == a.AppName &&
		c.CycleID == a.CycleID &&
		c.CycleName == a.CycleName &&
		c.Key == a.Key &&
		c.ServerAddr == a.ServerAddr
}

// ClientConfigs collection of client's config
type ClientConfigs struct {
	configList []*ClientConfig
	isCli      bool
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

// IsSetupRequired tell us whether in cli we need to setup the Client Config before
// procceding furthur
func (c *ClientConfigs) IsSetupRequired() bool {
	return len(c.configList) == 0
}

// load the config from file name with the given path
func (c *ClientConfigs) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(c.configList)
	if err != nil {
		return err
	}

	return nil
}

// save config into a file in given path
func (c *ClientConfigs) save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(c)
	if err != nil {
		return err
	}

	return nil
}

// Save the ClientConfig
func (c *ClientConfigs) Save() error {
	var err error

	if !c.isCli {
		return fmt.Errorf("WarpFile can't be saved in mobile environment")
	}

	err = c.save(iosWarpFilePath)
	if err != nil {
		return err
	}

	err = c.save(androidWarpFilePath)
	if err != nil {
		return err
	}

	return nil
}

// NewClientConfigsForCli creates a config for cli
func NewClientConfigsForCli() (*ClientConfigs, error) {
	iosConfigs := ClientConfigs{isCli: true}
	androidConfigs := ClientConfigs{isCli: true}

	err := iosConfigs.load(iosWarpFilePath)
	if err != nil {
		return nil, err
	}

	err = androidConfigs.load(androidWarpFilePath)
	if err != nil {
		return nil, err
	}

	if !isClientConfigsEqual(&iosConfigs, &androidConfigs) {
		return nil, fmt.Errorf("WarpFiles are mismatched")
	}

	// at this point, it doesn't really matter ios or android config
	// should be return, both are the same
	return &iosConfigs, nil
}

// NewClientConfigsForMobile creates a client config for mobile
func NewClientConfigsForMobile(bundlePath string) (*ClientConfigs, error) {
	bundlePathWarpFile := filepath.Join(bundlePath, "WarpFile")

	config := ClientConfigs{isCli: false}

	err := config.load(bundlePathWarpFile)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// isClientConfigsEqual compares two ClientConfigs. all the items must be in order in both of them
// order of config are matter
func isClientConfigsEqual(a, b *ClientConfigs) bool {
	if len(a.configList) != len(b.configList) {
		return false
	}

	for idx, config := range a.configList {
		if !config.equal(b.configList[idx]) {
			return false
		}
	}

	return true
}
