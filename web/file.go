package web

import (
	"io"
	"os"
)

// CopyDataToFile copies data from given input io.Reader and put it into
// a given file name.
func CopyDataToFile(in io.Reader, to string) error {
	out, err := os.Create(to)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)

	if err != nil {
		return err
	}

	return nil
}
