package cli

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"syscall"

	"strings"

	"net/url"

	"golang.org/x/crypto/ssh/terminal"
)

func Input(label string, isPassword bool) string {
	var value string

	fmt.Print(label)

	if isPassword {
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err == nil {
			value = string(bytePassword)
		}
		fmt.Println()
	} else {
		reader := bufio.NewReader(os.Stdin)
		value, _ = reader.ReadString('\n')
		value = strings.TrimSpace(value)
	}

	return value
}

func joinURL(base, target string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	u.Path = path.Join(u.Path, target)
	return u.String(), nil
}

func apiURL(path string) (string, error) {
	projectConfig, err := ProjectConfig()
	if err != nil {
		return "", err
	}

	url, err := joinURL(projectConfig.Server.Addr, path)
	if err != nil {
		return "", err
	}

	return url, nil
}
