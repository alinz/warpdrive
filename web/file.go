package web

import (
	"io"
	"os"
)

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
