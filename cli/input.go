package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"strconv"

	"golang.org/x/crypto/ssh/terminal"
)

func terminalInputExpect(label string, expect []string, defaultValue string) string {
	var cloneExpect []string
	for _, value := range expect {
		var target string
		if value == defaultValue {
			target = fmt.Sprintf("(%s)", value)
		} else {
			target = fmt.Sprintf("%s", value)
		}
		cloneExpect = append(cloneExpect, target)
	}

	label = fmt.Sprintf("%s %s ", label, cloneExpect)

	for {
		input := terminalInput(label, false)
		if input == "" {
			return defaultValue
		}

		for _, value := range expect {
			if value == input {
				return value
			}
		}
	}
}

func terminalInput(label string, isPassword bool) string {
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

func terminalInputAsInt64(label string, isPassword bool) (int64, error) {
	value := terminalInput(label, isPassword)
	return strconv.ParseInt(value, 10, 8)
}
