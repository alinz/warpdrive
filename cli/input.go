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
