package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
)

/**
 * Global Configuration
 */

//GlobalConfig store basic information about warpdrive
type GlobalConfig struct {
	// Sessions the key is the server url and the value of this key is jwt token
	// for authenticate to that server
	Sessions map[string]string `json:"sessions"`
}

// Path returns the global path which stores in user directory
// with the name of `.warp`
func (g *GlobalConfig) path() string {
	var home string
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	return path.Join(home, ".warp")
}

func (g *GlobalConfig) Save() error {
	file, err := os.Create(g.path())
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(g)
}

func (g *GlobalConfig) Load() error {
	file, err := os.Open(g.path())
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(g)
}

func (g *GlobalConfig) GetSessionFor(serverAddr string) (string, error) {
	session, ok := g.Sessions[serverAddr]
	if !ok {
		return "", fmt.Errorf("session not found for %s", serverAddr)
	}

	return session, nil
}

func (g *GlobalConfig) SetSessionFor(serverAddr, session string) {
	g.Sessions[serverAddr] = session
}

// NewGlobalConfig creates and load the global config
func NewGlobalConfig() *GlobalConfig {
	var config GlobalConfig
	config.Load()

	return &config
}
