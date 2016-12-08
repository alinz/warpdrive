package cli

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"encoding/json"

	"github.com/spf13/cobra"
)

type config interface {
	Load() error
	Save() error
}

/**
 * Global Configuration
 */

type globalConfig struct {
	// Sessions the key is the server url and the value of this key is jwt token
	// for authenticate to that server
	Sessions map[string]string `json:"sessions"`
}

// Path returns the global path which stores in user directory
// with the name of `.warp`
func (g *globalConfig) path() string {
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

func (g *globalConfig) Save() error {
	path := g.path()
	return configSave(path, g)
}

func (g *globalConfig) Load() error {
	path := g.path()
	return configLoad(path, g)
}

func (g *globalConfig) getSessionFor(serverAddr string) (string, error) {
	session, ok := g.Sessions[serverAddr]
	if !ok {
		return "", fmt.Errorf("session not found for %s", serverAddr)
	}

	return session, nil
}

func (g *globalConfig) setSessionFor(serverAddr, session string) {
	g.Sessions[serverAddr] = session
}

/**
 * Local Configuration
 */

type localConfig struct {
	ServerAddr string `json:"server_addr"`
	AppID      int64  `json:"app_id"`
	CycleID    int64  `json:"cycle_id"`
	Key        string `json:"key"`
}

func (l *localConfig) paths() []string {
	return []string{
		"android/app/src/main/assets/WarpFile",
		"ios/WarpFile",
	}
}

func (l *localConfig) Save() error {
	for _, loc := range l.paths() {
		// make sure the folder exists before creating a file.
		os.MkdirAll(path.Dir(loc), os.ModePerm)
		err := configSave(loc, l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *localConfig) Load() error {
	var err error
	loaded := false

	for _, path := range l.paths() {
		err = configLoad(path, l)
		if err == nil {
			loaded = true
			break
		}
	}

	if !loaded {
		return fmt.Errorf("'WarpFile' not found in '%s'", err.Error())
	}

	return nil
}

func (l *localConfig) isRequiredSetup() bool {
	return l.ServerAddr == "" || l.AppID == 0 || l.CycleID == 0
}

/**
 * Load and Save config files
 */

func configLoad(path string, conf config) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(path)
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(conf)
}

func configSave(path string, conf config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(conf)
}

/**
 * Root Command
 */

var RootCmd = &cobra.Command{
	Use:   "warp",
	Short: "In-App upgrade service for React-Native! Supporting iOS and Android apps",
	Long: `
A Fast and Flexible upgrade service for React-Native apps!
loved by alinz and Pressly Inc.

Complete documentation is available at https://pressly.github.io/warpdrive
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("please run 'warp -h' for usage")
	},
}

func init() {}
