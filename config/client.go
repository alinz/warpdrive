package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AppConfig represents app basic information in config
type AppConfig struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (a *AppConfig) equal(b *AppConfig) bool {
	return a.ID == b.ID && a.Name == b.Name
}

//CycleConfig data represents each cycle in config
type CycleConfig struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (c *CycleConfig) equal(b *CycleConfig) bool {
	return c.ID == b.ID &&
		c.Name == b.Name &&
		c.Key == b.Key
}

// ClientConfig is being used in both cli and client bundle in mobile app
type ClientConfig struct {
	ServerAddr string         `json:"server_addr"`
	App        AppConfig      `json:"app"`
	Cycles     []*CycleConfig `json:"cycles"`
	isCli      bool
	paths      []string
}

func (c *ClientConfig) equal(b *ClientConfig) bool {
	if c.ServerAddr != b.ServerAddr ||
		!c.App.equal(&b.App) {
		return false
	}

	if len(c.Cycles) != len(b.Cycles) {
		return false
	}

	for idx, cycle := range c.Cycles {
		if !cycle.equal(b.Cycles[idx]) {
			return false
		}
	}

	return true
}

// Load the config from file name with the given path
func (c *ClientConfig) Load() error {
	var err error
	var prev *ClientConfig

	if len(c.paths) == 0 {
		return fmt.Errorf("paths not set")
	}

	for _, path := range c.paths {
		err = func(path string) error {
			curr := &ClientConfig{}

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()

			err = json.NewDecoder(file).Decode(curr)
			if err != nil {
				return err
			}

			if prev != nil && !prev.equal(curr) {
				return fmt.Errorf("WarpFiles not matched")
			}

			prev = curr
			return nil
		}(path)

		if err != nil {
			return err
		}
	}

	c.App = prev.App
	c.Cycles = prev.Cycles
	c.ServerAddr = prev.ServerAddr

	return nil
}

// Save config into a file in given path
func (c *ClientConfig) Save() error {
	if !c.isCli {
		return fmt.Errorf("mobile can't save WarpFile")
	}

	if len(c.paths) == 0 {
		return fmt.Errorf("paths not set")
	}

	for _, path := range c.paths {
		err := func(path string) error {
			file, err := os.Create(path)
			if err != nil {
				return err
			}

			defer file.Close()

			return json.NewEncoder(file).Encode(c)
		}(path)

		if err != nil {
			return err
		}
	}

	return nil
}

// IsSetupRequired will let us know whether the config is empty or not
func (c *ClientConfig) IsSetupRequired() bool {
	return len(c.Cycles) == 0
}

// AddCycle adds a new cycle into the list, prevent duplicates, order are also matter
func (c *ClientConfig) AddCycle(newCycleConfig *CycleConfig) error {
	cycleConfig, err := c.GetCycle(newCycleConfig.Name)
	if err == nil {
		if cycleConfig.ID == newCycleConfig.ID {
			return fmt.Errorf("cycle '%s' already exists", newCycleConfig.Name)
		}
		return fmt.Errorf("cycle '%s' has different id", newCycleConfig.Name)
	}

	c.Cycles = append(c.Cycles, newCycleConfig)
	return nil
}

// GetCycle gets config based on app's and cycle's names.
func (c *ClientConfig) GetCycle(cycleName string) (*CycleConfig, error) {
	for _, cycle := range c.Cycles {
		if cycle.Name == cycleName {
			return cycle, nil
		}
	}

	return nil, fmt.Errorf("config not found")
}

// NewClientConfigsForCli creates a config for cli
func NewClientConfigsForCli() *ClientConfig {

	// at this point, it doesn't really matter ios or android config
	// should be return, both are the same
	return &ClientConfig{
		isCli: true,
		paths: []string{
			"ios/WarpFile",
			"android/app/src/main/assets/WarpFile",
		},
	}
}

// NewClientConfigsForMobile creates a client config for mobile
func NewClientConfigsForMobile(bundlePath string) *ClientConfig {
	bundlePathWarpFile := filepath.Join(bundlePath, "WarpFile")
	return &ClientConfig{
		isCli: false,
		paths: []string{
			bundlePathWarpFile,
		},
	}
}
